# Installation Guide

## Import the `cgo.dump` File

### Step1. Install Docker

If Docker is not already installed, follow these steps:

1. **Linux/MacOS**:
   - Refer to the official documentation to install Docker: https://docs.docker.com/get-docker/
2. **Windows**:
   - Install Docker Desktop: https://www.docker.com/products/docker-desktop/

After installation, verify Docker is running with the following command:

```bash
docker --version
```

------

### Step2. Install PostgreSQL 13 Using Docker

1. Pull the PostgreSQL 13 image:

   ```bash
   docker pull postgres:13
   ```

2. Create and run a PostgreSQL container:

   ```bash
   docker run -d \
     --name postgres13 \
     -e POSTGRES_USER=your_username \
     -e POSTGRES_PASSWORD=your_password \
     -e POSTGRES_DB=your_database \
     -p 5432:5432 \
     postgres:13
   ```

   - **POSTGRES_USER**: Set the PostgreSQL username (e.g., `postgres`).
   - **POSTGRES_PASSWORD**: Set the PostgreSQL password.
   - **POSTGRES_DB**: Set the default database name.
   - **-p 5432:5432**: Map port 5432 on your machine to port 5432 in the container.

3. Check the container status:

   ```bash
   docker ps
   ```
   If successful, you will see a container named `postgres13`.

4. Access the PostgreSQL container:

   ```bash
   docker exec -it postgres13 bash
   ```

### Step3. Copy the `cgo.dump` File into the Container

Run the following command on your host machine to copy `cgo.dump` into the container:

```bash
docker cp cgo.dump postgres13:/tmp/cgo.dump
```

### Step4. Import the `cgo.dump` File

1. Access the container and switch to the PostgreSQL user:

   ```bash
   docker exec -it postgres13 bash
   su - postgres
   ```

2. Restore the database using `pg_restore`:

   ```bash
   pg_restore -U your_username -d your_database /tmp/cgo.dump
   ```

   - **your_username**: Specify the database username.
   - **your_database**: Specify the target database name.

3. To clean the target database before importing, add the `--clean` option:

   ```bash
   pg_restore -U your_username -d your_database --clean /tmp/cgo.dump
   ```

4. Verify the import: Log into the database and check the data:

   ```bash
   psql -U your_username -d your_database
   \dt
   ```

### Step5(Optional). Access the Imported Data using GUI

You can install PgAdmin to access the imported data using a graphical interface. You can download the installer from the [official website](https://www.pgadmin.org/download/)

## Build our Modified Go Toolchain based on Go1.17.7

### Step1. Download Go1.7.7 toolchain source code

1. download Go1.17.7 toolchain source code: https://go.dev/dl/go1.17.7.src.tar.gz
```bash
wget https://go.dev/dl/go1.17.7.src.tar.gz
```
2. extract toolchain source code

```bash
tar -xzvf go1.17.7.src.tar.gz
```

### Step2. Use our `go1.17.7_cgoptr.tar.gz`

1. extract our `go1.17.7_cgoptr.tar.gz`
```bash
tar -xzvf go1.17.7_cgoptr.tar.gz
```
2. overwrite the original Go1.17.7 toolchain source with the extracted file
```bash
cd go1.17.7_cgoptr
cp -r * ../go/
```
3. build the modifiled toochain

```bash
cd ../go/src
./make.bash # build
```

note: to build the modified toolchain, your machine must have Go installation (need the Go compiler to build the Go toolchain). ref [Download and install - The Go Programming Language](https://go.dev/doc/install) for details.

4. verify your build

```bash
cd ../bin
./go --version
```

you can add this to your `$PATH`. and modify your `$GOROOT` if you want to use the modified toolchain later. 

see [Go Wiki: InstallTroubleshooting - The Go Programming Language](https://go.dev/wiki/InstallTroubleshooting) for how to set `$GOROOT`.
