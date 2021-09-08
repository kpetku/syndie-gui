module github.com/kpetku/syndie-gui

go 1.16

replace github.com/kpetku/syndie-core => /home/desktop/go/src/github.com/kpetku/syndie-core

replace github.com/kpetku/libsyndie => /home/desktop/go/src/github.com/kpetku/libsyndie

require (
	fyne.io/fyne/v2 v2.0.4
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-gl/gl v0.0.0-20210905235341-f7a045908259 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20210727001814-0db043d8d5be // indirect
	github.com/kpetku/libsyndie v1.0.0
	github.com/kpetku/syndie-core v1.0.1
	github.com/srwiley/oksvg v0.0.0-20210519022825-9fc0c575d5fe // indirect
	github.com/srwiley/rasterx v0.0.0-20210519020934-456a8d69b780 // indirect
	go.etcd.io/bbolt v1.3.6
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/net v0.0.0-20210907225631-ff17edfbf26d
	golang.org/x/sys v0.0.0-20210906170528-6f6e22806c34 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
