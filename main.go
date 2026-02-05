package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CloudWatchAlarms []string `yaml:"cloudwatch_alarm_list"`
}

type CloudWatchManager struct {
	client *cloudwatch.Client
	ctx    context.Context
}

var (
	environment = flag.String("env", "currentsite-dev", "Environment name (e.g., currentsite-dev)")
	region      = flag.String("region", "ap-northeast-1", "AWS region (default: ap-northeast-1)")
)

// NewCloudWatchManager creates a new CloudWatchManager with the provided AWS configuration
func NewCloudWatchManager(ctx context.Context) (*CloudWatchManager, error) {
	cfg, err := loadAWSConfig(ctx, *environment, *region)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return &CloudWatchManager{
		client: cloudwatch.NewFromConfig(cfg),
		ctx:    ctx,
	}, nil
}

// AlarmExists checks if a cloudwatch alarm with the alarm name exists
func (cw *CloudWatchManager) AlarmExists(alarmName string) (bool, error) {
	input := &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{alarmName},
	}

	result, err := cw.client.DescribeAlarms(cw.ctx, input)
	fmt.Println("DescribeAlarms result:", result)
	if err != nil {
		return false, fmt.Errorf("error describing alarm %s: %w", alarmName, err)
	}

	// Check if alarm exists in MetricAlarms or CompositeAlarms
	if len(result.MetricAlarms) > 0 || len(result.CompositeAlarms) > 0 {
		return true, nil
	}

	return false, nil
}

// DeleteAlarms deletes the specified cloudwatch alarms
func (cw *CloudWatchManager) DeleteAlarm(alarmName string) error {
	input := &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []string{alarmName},
	}

	_, err := cw.client.DeleteAlarms(cw.ctx, input)
	if err != nil {
		return fmt.Errorf("error deleting alarm %s: %w", alarmName, err)
	}

	return nil
}

// LoadAlarmsFromYAML loads the list of CloudWatch alarms from a YAML file.
func LoadAlarmsFromYAML(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return config.CloudWatchAlarms, nil
}

func main() {
	flag.Parse()

	ctx := context.Background()

	// Load alarms from yaml file
	fmt.Println("Loading CloudWatch alarms from YAML file...")
	alarmNames, err := LoadAlarmsFromYAML("cloudwatch-alarms.yaml")
	if err != nil {
		fmt.Printf("Error loading alarms: %v\n", err)
		return
	}
	fmt.Printf("Found %d CloudWatch alarms\n", len(alarmNames))

	cwManager, err := NewCloudWatchManager(ctx)
	if err != nil {
		log.Fatalf("Failed to create CloudWatch manager: %v", err)
	}

	// Track statistics
	var existingAlarms, nonExistentAlarms, deletedAlarms, failedDeletions []string

	// Step 1: Check which alarms exists
	fmt.Println("====== Checking alarms Existence ======")

	for i, alarmName := range alarmNames {
		fmt.Printf("[%d/%d] Checking alarm: %s... ", i+1, len(alarmNames), alarmName)
		exists, err := cwManager.AlarmExists(alarmName)
		if err != nil {
			fmt.Printf("Error checking alarm: %v\n", err)
			continue
		}
		if exists {
			existingAlarms = append(existingAlarms, alarmName)
			fmt.Println("Exists")
		} else {
			nonExistentAlarms = append(nonExistentAlarms, alarmName)
			fmt.Println("Does not exist")
		}
	}

	// Step 2: Confirm deletion
	if len(existingAlarms) == 0 {
		fmt.Println("\nNo alarms to delete. Exiting.")
		return
	}

	fmt.Printf("\n=== Deletion Confirmation ===\n")
	fmt.Printf("About to delete %d alarm(s). Do you want to continue? (yes/no): ", len(existingAlarms))

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "yes" && confirmation != "y" && confirmation != "YES" && confirmation != "Y" {
		fmt.Println("Deletion cancelled by user.")
		return
	}

	// Step 3: Delete existing alarms
	fmt.Println("\n====== Deleting Alarms ======")

	for i, alarmName := range existingAlarms {
		fmt.Printf("[%d/%d] Deleting alarm: %s... ", i+1, len(existingAlarms), alarmName)

		err := cwManager.DeleteAlarm(alarmName)
		if err != nil {
			log.Printf("Failed to delete: %v", err)
			failedDeletions = append(failedDeletions, alarmName)
			fmt.Println("FAILED ✗")
			continue
		}

		deletedAlarms = append(deletedAlarms, alarmName)
		fmt.Println("DELETED ✓")
	}

	// Final summary
	fmt.Printf("\n=== Final Summary ===\n")
	fmt.Printf("Successfully deleted: %d alarm(s)\n", len(deletedAlarms))
	fmt.Printf("Failed deletions: %d alarm(s)\n", len(failedDeletions))

	if len(failedDeletions) > 0 {
		fmt.Println("\nFailed to delete:")
		for _, name := range failedDeletions {
			fmt.Printf("  - %s\n", name)
		}
	}

	if len(deletedAlarms) > 0 {
		fmt.Println("\nSuccessfully deleted:")
		for _, name := range deletedAlarms {
			fmt.Printf("  - %s\n", name)
		}
	}

	fmt.Println("\nOperation completed!")
}

// loadAWSConfig loads AWS configuration with the specified environment and region.
func loadAWSConfig(ctx context.Context, env, region string) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	switch env {
	case "currentsite-dev":
		opts = append(opts, config.WithSharedConfigProfile("currentsite-dev"))
	case "currentsite-prod":
		opts = append(opts, config.WithSharedConfigProfile("currentsite-prod"))
	default:
		// if no env is provided, default to currentsite-dev
		opts = append(opts, config.WithSharedConfigProfile("currentsite-dev"))
	}

	// Set region if provided
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	} else {
		opts = append(opts, config.WithRegion("ap-northeast-1"))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}
