# go_miner
go_miner is a SQL database miner written in golang. 

**Author**: [th3jiv3r][twitter]

### New Features!
  - map MS-SQL server databases
  - uses regex to look for specified keywords in table columns, showing if there are any matches

### Installation
#### Setup and build go_miner
```sh
$ go get github.com/goodlandsecurity/go_miner/dbminer
$ git clone https://github.com/goodlandsecurity/go_miner
$ cd go_miner
$ go build go_miner.go
```
#### Setup test MS-SQL docker instance
Use the transactions_seed.txt file to build the application database if you want to have two databases configured
```sh
$ docker run --name mssql -p 1433:1433 -e 'ACCEPT_EULA=Y' \
-e 'SA_PASSWORD=P@$$w0rd!' -d microsoft/mssql-server-linux
$ docker exec -it mssql /opt/mssql-tools/bin/sqlcmd -S localhost \
-U sa -P 'P@$$w0rd!'
> CREATE DATABASE store;
> GO
> USE store;
> CREATE TABLE transactions(ccnum varchar(32), date date, amount decimal(7,2),cvv char(4), exp date);
> GO
> INSERT INTO transactions(ccnum, date, amount, cvv, exp) VALUES('4444333322221111', '2019-01-05', 100.12, '1234', '2020-09-01');
> INSERT INTO transactions(ccnum, date, amount, cvv, exp) VALUES('4444123456789012', '2019-01-07', 2400.18, '5544', '2021-02-01');
> INSERT INTO transactions(ccnum, date, amount, cvv, exp) VALUES('4465122334455667', '2019-01-29', 1450.87, '9876', '2020-06-01');
> GO
```
### Example Use:  
  - *go_miner -host localhost -user sa -password P@$$w0rd!*
  - *go_miner -host 10.0.0.1 -port 1234 -user sa -password P@$$w0rd!*
  - *go_miner -host 192.168.0.1 -user sa -password P@$$w0rd! -debug true* 

#### License
  - GNU General Public License v3.0


[twitter]: <https://twitter.com/th3_jiv3r>
