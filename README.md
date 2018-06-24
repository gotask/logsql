# logsql
you can search any log by sql command using logsql.
logsql is based on TextQL (https://github.com/dinedal/textql), so it has the same License with TextQL.
use like this:
```
logsql -s "select * from a" -f "|" a.log
//out:1 2 3		a.log:1|2|3

logsql -s "select x from a" -h a.csv
//out:1		
a.csv:
x y z
1 2 3
	  
logsql -s "select x+y from a" -h a.xlsx
//out:3		
a.xlsx:
x y z
1 2 3
```
you can use it without "select" or "from":
```
logsql -s "*" -f "|" a.log
logsql -s "x" -h a.csv
```
you can see TextQL for anyother complex usage.
by the way,it reads xlsx by https://github.com/tealeg/xlsx.

### download the executable file in release(windows or linux).