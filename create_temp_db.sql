CREATE TABLE IF NOT EXISTS temp_and_humidity (
  timestamp DATETIME, 
  temp_degrees_c REAL, 
  relative_humidity REAL, 
  sensor TEXT, 
  debug TEXT
);
