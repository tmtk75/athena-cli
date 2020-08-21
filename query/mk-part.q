ALTER TABLE cloudtrail_logs ADD 
   PARTITION (region='{{ .region }}',
              year='{{ .year }}',
              month='{{ .month }}',
              day='{{ .day }}')
   LOCATION 's3://{{ .s3bucketname }}/AWSLogs/{{ .accountId }}/CloudTrail/{{ .region }}/{{ .year }}/{{ .month }}/{{ .day }}/';

