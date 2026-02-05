# CloudWatch Alarm Deletion Tool

This Go program reads CloudWatch alarm names from a YAML file, checks if they exist, and deletes them after confirmation.

## Features

- ✅ Reads alarm names from YAML file
- ✅ Checks if each alarm exists in AWS CloudWatch
- ✅ Shows alarm details (state, metric name)
- ✅ Confirms before deletion
- ✅ Deletes alarms in a loop
- ✅ Comprehensive error handling and reporting

## Prerequisites

1. Go 1.21 or higher
2. AWS credentials configured (via `~/.aws/credentials` or environment variables)
3. Appropriate IAM permissions for CloudWatch:
   - `cloudwatch:DescribeAlarms`
   - `cloudwatch:DeleteAlarms`

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Configure AWS credentials:
```bash
saml2aws login -a currentsite-{dev|prod} --profile currentsite-{dev|prod} 
```

3. Update `cloudwatch-alarms.yml` with your alarm names:
```yaml
cloudwatch_alarms:
  - my-alarm-1
  - my-alarm-2
  - my-alarm-3
```

## Usage

Run the program:
```bash
go run main.go
```

Binary:
```Bash
go build

chmod +x cloudwatch-alarm-delete

./chmod +x cloudwatch-alarm-delete 
```

The program will:
1. Load alarm names from `cloudwatch-alarms.yml`
2. Check which alarms exist in CloudWatch
3. Show a summary of existing vs non-existent alarms
4. Ask for confirmation before deletion
5. Delete the confirmed alarms
6. Show final results

## Example Output

```
Loading alarms from cloudwatch-alarms.yml...
Found 31 alarms in YAML file

=== Checking Alarm Existence ===
[1/31] Checking alarm: my-alarm-1... EXISTS ✓
    State: OK, Metric: CPUUtilization
[2/31] Checking alarm: my-alarm-2... NOT FOUND ✗
[3/31] Checking alarm: my-alarm-3... EXISTS ✓
    State: ALARM, Metric: NetworkIn

=== Existence Check Summary ===
Total alarms in YAML: 31
Existing alarms: 25
Non-existent alarms: 6

Non-existent alarms:
  - my-alarm-2
  - my-alarm-7
  ...

=== Deletion Confirmation ===
About to delete 25 alarm(s). Do you want to continue? (yes/no): yes

=== Deleting Alarms ===
[1/25] Deleting alarm: my-alarm-1... DELETED ✓
[2/25] Deleting alarm: my-alarm-3... DELETED ✓
...

=== Final Summary ===
Successfully deleted: 25 alarm(s)
Failed deletions: 0 alarm(s)

Operation completed!
```

## YAML File Format

```yaml
cloudwatch_alarms:
  - alarm-name-1
  - alarm-name-2
  - alarm-name-3
  # Add more alarms as needed
```

```Bash
=== Deletion Confirmation ===
About to delete 31 alarm(s). Do you want to continue? (yes/no): yes

====== Deleting Alarms ======
[1/31] Deleting alarm: api-test-activity-cs... DELETED ✓
[2/31] Deleting alarm: awsec2-rd-group-CPU-Utilization... DELETED ✓
[3/31] Deleting alarm: awsec2-rd-group-High-CPU-Utilization... DELETED ✓
[4/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_additional_info... DELETED ✓
[5/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_additional_info-dev... DELETED ✓
[6/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_booking_info... DELETED ✓
[7/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_booking_info-dev... DELETED ✓
[8/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_country... DELETED ✓
[9/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_country-dev... DELETED ✓
[10/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_currency_rate... DELETED ✓
[11/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_currency_rate-dev... DELETED ✓
[12/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_hotel... DELETED ✓
[13/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_hotel-dev... DELETED ✓
[14/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_origin... DELETED ✓
[15/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_origin-dev... DELETED ✓
[16/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_origin_timezone... DELETED ✓
[17/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_origin_timezone-dev... DELETED ✓
[18/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_unit... DELETED ✓
[19/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_m_unit-dev... DELETED ✓
[20/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_navi_category... DELETED ✓
[21/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_navi_category-dev... DELETED ✓
[22/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_navi_origin... DELETED ✓
[23/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_navi_origin-dev... DELETED ✓
[24/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_ptr_basic_info... DELETED ✓
[25/31] Deleting alarm: dev_master_api_error_alarm_lambda_get_ptr_basic_info-dev... DELETED ✓
[26/31] Deleting alarm: dev_master_api_error_alarm_lambda_timeout_handler... DELETED ✓
[27/31] Deleting alarm: dev_master_api_error_alarm_lambda_tr_nologin_lambda_timeout_handler... DELETED ✓
[28/31] Deleting alarm: dev_master_api_error_alarm_lambda_tr_nologin_lambda_timeout_handler-dev... DELETED ✓
[29/31] Deleting alarm: dev_partner_api_error_alarm_lambda_get_ptr_basic_info... DELETED ✓
[30/31] Deleting alarm: dev_partner_api_error_alarm_lambda_partner_timeout_handler... DELETED ✓
[31/31] Deleting alarm: rds_connection_high... DELETED ✓

=== Final Summary ===
Successfully deleted: 31 alarm(s)
Failed deletions: 0 alarm(s)

Successfully deleted:
  - api-test-activity-cs
  - awsec2-rd-group-CPU-Utilization
  - awsec2-rd-group-High-CPU-Utilization
  - dev_master_api_error_alarm_lambda_get_m_additional_info
  - dev_master_api_error_alarm_lambda_get_m_additional_info-dev
  - dev_master_api_error_alarm_lambda_get_m_booking_info
  - dev_master_api_error_alarm_lambda_get_m_booking_info-dev
  - dev_master_api_error_alarm_lambda_get_m_country
  - dev_master_api_error_alarm_lambda_get_m_country-dev
  - dev_master_api_error_alarm_lambda_get_m_currency_rate
  - dev_master_api_error_alarm_lambda_get_m_currency_rate-dev
  - dev_master_api_error_alarm_lambda_get_m_hotel
  - dev_master_api_error_alarm_lambda_get_m_hotel-dev
  - dev_master_api_error_alarm_lambda_get_m_origin
  - dev_master_api_error_alarm_lambda_get_m_origin-dev
  - dev_master_api_error_alarm_lambda_get_m_origin_timezone
  - dev_master_api_error_alarm_lambda_get_m_origin_timezone-dev
  - dev_master_api_error_alarm_lambda_get_m_unit
  - dev_master_api_error_alarm_lambda_get_m_unit-dev
  - dev_master_api_error_alarm_lambda_get_navi_category
  - dev_master_api_error_alarm_lambda_get_navi_category-dev
  - dev_master_api_error_alarm_lambda_get_navi_origin
  - dev_master_api_error_alarm_lambda_get_navi_origin-dev
  - dev_master_api_error_alarm_lambda_get_ptr_basic_info
  - dev_master_api_error_alarm_lambda_get_ptr_basic_info-dev
  - dev_master_api_error_alarm_lambda_timeout_handler
  - dev_master_api_error_alarm_lambda_tr_nologin_lambda_timeout_handler
  - dev_master_api_error_alarm_lambda_tr_nologin_lambda_timeout_handler-dev
  - dev_partner_api_error_alarm_lambda_get_ptr_basic_info
  - dev_partner_api_error_alarm_lambda_partner_timeout_handler
  - rds_connection_high

Operation completed!
jabes.pauya@A808 cloudwatch-alarms-delete
```

## Error Handling

The program handles various scenarios:
- Missing or invalid YAML file
- AWS authentication errors
- Alarms that don't exist
- Deletion failures (e.g., permission issues)
- Network errors

## Safety Features

- ✅ Confirmation prompt before deletion
- ✅ Shows which alarms exist before attempting deletion
- ✅ Only deletes existing alarms
- ✅ Detailed logging of success/failure for each alarm
- ✅ Final summary report

## Build

To build an executable:
```bash
go build -o cloudwatch-alarm-deleter
./cloudwatch-alarm-deleter
```

## License

MIT