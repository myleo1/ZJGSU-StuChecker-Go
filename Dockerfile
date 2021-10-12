FROM alpine:latest
LABEL maintainer="leo <leo@leom.me>" \
	version="v0.0.1-beta" \
	description="ZJGSU-StuChecker-Go"
WORKDIR /root
ADD build/ZJGSU-StuChecker-Go init.sh /root/
RUN echo '04       23      *       *       *       /root/ZJGSU-StuChecker-Go' > /etc/crontabs/root \
	&& chmod -R 777 /root
CMD ["/bin/sh","/root/init.sh"]