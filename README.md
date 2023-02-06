# ACMHomepage

ACMHomepage is a Homepage for ACMer.

## Build and Run

This section is used to build and run the service. Not for testing.

The only requirement is Docker.

1. The first thing you should do is build the compose:

   ```shell
   docker compose build
   ```

   If Go cannot download the dependencies modules, you can try setting the
   `GOPROXY` build argument:

   ```shell
   docker compose build --build-arg GOPROXY=<proxy-url>
   ```

   You can find useful information on the Internet regarding setting a proxy for
   Go.

2. Run the service by using the following command:

   ```
   docker compose up
   ```

