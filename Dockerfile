FROM centos:centos7.9.2009 AS Builder

WORKDIR /root

#copying any kind of dependencies
COPY ./cmd/api/compile/dependencies .
#copying uploaded files
COPY ./cmd/api/compile/config-file-downloads/ .
#copying release directory, where all the returning files from the container will go to
COPY ./cmd/api/compile/release .
#unpack any necessary dependencies
RUN // //

#give permission to execute any automation scripts and run them
RUN chmod 777 ./*.sh &&\
    chmod 777 ./*.exp &&\
    ./update.sh
#run any scripts to monitor the given instructions
CMD ["./def-watch.sh", "&"]

#copy the compiled file to a smaller container
#FROM alpine:latest
#WORKDIR /root
#COPY --from=builder /root/CTCC_SGU_V2.0.2B03_210720_B/compatible_branch/make/target/bin/1GPON/release/P24_HGUV6.5.3B01_UPGRADE.bin ./release/
