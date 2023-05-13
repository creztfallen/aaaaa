FROM centos:centos7.9.2009 AS Builder

WORKDIR /root

#copying any kind of dependencies
COPY ./cmd/api/compile/dependencies .
#copying downloaded filessudo
COPY ./cmd/api/compile/config-file-downloads/ .

COPY ./cmd/api/compile/release .

RUN tar -xzvf CTCC_SGU_V2.0.2B03_210720_B_20221129.tar.gz

#give permission to execute any automation scripts and run them
RUN chmod 777 ./*.sh &&\
    chmod 777 ./*.exp &&\
    ./update.sh

CMD ["./def-watch.sh", "&"]

#copy the compiled file to a smaller container
#FROM alpine:latest
#WORKDIR /root
#COPY --from=builder /root/CTCC_SGU_V2.0.2B03_210720_B/compatible_branch/make/target/bin/1GPON/release/P24_HGUV6.5.3B01_UPGRADE.bin ./release/

#sudo docker run --name sentinel -v /home/mmoreno/Documents/projects/compiler_api/cmd/api/compile/release:/root/CTCC_SGU_V2.0.2B03_210720_B/compatible_branch/make/target/bin/1GPON/release/ -dit creztfallen/compiler-environment bash
#sudo docker cp ./def-decry-cfg sentinel:/root/
