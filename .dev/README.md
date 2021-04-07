## Development Instruction
- copy Vagrantfile and setup.sh to custom folder
```
cd .dev
cp Vagrantfile ~/<custom folder>/Vagrantfile

cp setuo.sh ~/<custom folder>/setup.sh
```
### start vagrant box
```# vagrant up```

```# vagrant ssh```

### remote debug

Configure remove debug  on intellij on port : 2345
- make sure project is cloned to vagrant shared folder
- build executable for debug on vagrant box
```
make build_debug
```
- start dlv
```
make dlv
```
- start remote debug from intellij
