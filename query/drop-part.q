ALTER TABLE cloudtrail_logs
 DROP PARTITION (year='{{ .year }}',
                 month='{{ .month }}',
		 day='{{ .day }}');

