package plugin

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// SetupTimeout 每个插件初始化最长超时时间，如果某个插件确实需要加载很久，可以自己修改这里的值
	SetupTimeout = 3 * time.Second

	// MaxPluginSize  最大插件个数
	MaxPluginSize = 1000
)

var (
	ErrDependsOnItself = fmt.Errorf("plugin not allowed to depend on itself")
)

var (
	factories = make(map[string]Factory) // plugin type => { plugin name => plugin factory }
	done      = make(chan struct{})
) // 插件初始化完成通知channel

type Depended interface {
	// DependsOn 假如一个插件依赖另一个插件，则返回被依赖的插件的列表：数组元素为 type-name 如 [ "selector-polaris" ]
	DependsOn() []string
}

// FlexDepended 弱依赖接口，如果被依赖的插件存在，才去保证被依赖的插件先初始化完成
type FlexDepended interface {
	FlexDependsOn() []string
}

type Info struct {
	factory Factory
	name    string
	cfg     yaml.Node
}

// Setup 初始化单个插件。
func (p *Info) Setup() error {
	var (
		ch  = make(chan struct{})
		err error
	)
	go func() {
		err = p.factory.Setup(p.name, &YamlNodeDecoder{Node: &p.cfg})
		close(ch)
	}()
	select {
	case <-ch:
	case <-time.After(SetupTimeout):
		return fmt.Errorf("setup plugin %s timeout", p.Key())
	}
	if err != nil {
		return fmt.Errorf("setup plugin %s error: %v", p.Key(), err)
	}
	return nil
}

// Depends 判断是否有依赖的插件未初始化过。
// 输入参数为所有插件的初始化状态。
// 输出参数bool true被依赖的插件未初始化完成，仍有依赖，false没有依赖其他插件或者被依赖的插件已经初始化完成
func (p *Info) Depends(setupStatus map[string]bool) (bool, error) {
	deps, ok := p.factory.(Depended)
	if !ok { // 该插件不依赖任何其他插件
		return false, nil
	}
	depends := deps.DependsOn()
	for _, name := range depends {
		if name == p.Key() {
			return false, ErrDependsOnItself
		}
		setup, ok := setupStatus[name]
		if !ok {
			return false, fmt.Errorf("depends plugin %s not exists", name)
		}
		if !setup {
			return true, nil
		}
	}
	return false, nil
}

// flexDepends 弱依赖，类似 Depends 方法，判断是否有依赖的插件存在且未初始化
func (p *Info) flexDepends(setupStatus map[string]bool) (bool, error) {
	fd, ok := p.factory.(FlexDepended)
	if !ok { // 不存在弱依赖关系
		return false, nil
	}
	depends := fd.FlexDependsOn()
	for _, name := range depends {
		if name == p.Key() {
			return false, ErrDependsOnItself
		}
		setup, ok := setupStatus[name]
		if !ok {
			return false, nil
		}
		if !setup {
			return true, nil
		}
	}
	return false, nil
}

// Key 插件的唯一索引：type-name 。
func (p *Info) Key() string {
	return fmt.Sprintf("%s", p.name)
}

// Config 插件统一配置 plugin type => { plugin name => plugin config } 。
type Config map[string]yaml.Node

// Setup 通过配置生成并装载具体插件。
func (c Config) Setup() error {
	var (
		pluginChan  = make(chan Info, MaxPluginSize) // 初始化插件队列，方便后面按顺序逐个加载插件
		setupStatus = make(map[string]bool)          // 插件初始化状态，plugin key => true初始化完成 false未初始化
	)

	// 从框架配置文件中逐个取出插件并放进channel队列中
	for name, cfg := range c {
		factory := Get(name)
		if factory == nil {
			return fmt.Errorf("plugin %s no registered or imported, do not configure", name)
		}
		p := Info{
			factory: factory,
			name:    name,
			cfg:     cfg,
		}
		select {
		case pluginChan <- p:
		default:
			return fmt.Errorf("plugin number exceed max limit:%d", len(pluginChan))
		}
		setupStatus[p.Key()] = false
	}

	// 从channel队列中取出插件并初始化
	num := len(pluginChan)
	for num > 0 {
		for i := 0; i < num; i++ {
			p := <-pluginChan
			if err := p.Setup(); err != nil {
				return err
			}
			setupStatus[p.Key()] = true
		}
		if len(pluginChan) == num { // 循环依赖导致无插件可以初始化，返回失败
			return fmt.Errorf("cycle depends, not plugin is setup")
		}
		num = len(pluginChan)
	}

	// 发出插件初始化完成通知，个别业务逻辑需要依赖插件完成才能继续往下执行
	select {
	case <-done: // 已经close过了
	default:
		close(done)
	}
	return nil
}

type YamlNodeDecoder struct {
	Node *yaml.Node
}

func (d *YamlNodeDecoder) Decode(cfg interface{}) error {
	{
		if d.Node == nil {
			return fmt.Errorf("yaml node empty")
		}
		return d.Node.Decode(cfg)
	}
}
