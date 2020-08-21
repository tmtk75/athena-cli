select max(eventTime) as maxEventTime from cloudtrail_logs
  where
    useridentity.principalid like '%{{ .principal }}%'
    and year = '2019' and month = '08' and day = '19'
