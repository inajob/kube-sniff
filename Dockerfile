FROM scratch

COPY ./kube-sniff /kube-sniff
COPY ./static /static
EXPOSE 9999

CMD ["/kube-sniff"]