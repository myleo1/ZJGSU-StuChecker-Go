FROM alpine:latest
LABEL maintainer="leo <admin@leow.tech>" \
	version="v1.0.0-beta" \
	description="ZJGSU-StuChecker-Go"
WORKDIR /root
ADD build/ZJGSU-StuChecker-Go /root/
RUN echo '04       23      *       *       *       /root/ZJGSU-StuChecker-Go yzy' > /etc/crontabs/root \
	&& chmod +x /root/ZJGSU-StuChecker-Go
CMD ["crond","-f"]