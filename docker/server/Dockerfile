FROM alpine
ADD server /bin
ADD newdag.cfg /
ADD peers.json /
ADD newdag_chain.db /
EXPOSE 1337 1338 80
ENV HOME=/
ENTRYPOINT ["server"]
CMD [] 
