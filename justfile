#* list cmd
h:
    just -l

#* inserting log
r:
    go run . > main.log

#* parsing log
#? tail -f = tail but still reading new changes
#? grep --line-buffered = ouput perline not bulk line
#? sed -u = unbuffered (send result as fast as possible)
p:
    tail -f main.log \
        | grep --line-buffered '^logdy ' \
        | sed -u 's/logdy //' \
        | logdy --config logdy.json