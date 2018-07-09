FROM alpine
ADD client /bin
ADD index.html /
ADD newdag.cfg /
ADD newdag_chain.db /
ADD newdag_wallet.dat /
EXPOSE 1339
EXPOSE 1388
ENV HOME=/
ENTRYPOINT ["client"]
CMD [] 
