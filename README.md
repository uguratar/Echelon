#Echelon For Seyes

Simple web service for writing logs to file for transporting to Logstash/Kibana.
Needs .htpasswd file at the same directory with compiled binary to run.

Paremeters:
 - Listening Port (eg: 3000)
 - Absolute path to write log files (eg: /var/log/my-logs/)


 Thanks to: 
 - https://bitbucket.org/kardianos/osext
 - https://github.com/abbot/go-http-auth
 - https://github.com/alecthomas/log4go
 - https://github.com/scorredoira/email
