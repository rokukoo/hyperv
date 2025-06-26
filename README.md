# hyperv
 
## 简介

hyperv 是一个用于管理 Microsoft Hyper-V 虚拟化环境的 Go 语言 SDK，基于 WMI 实现，支持虚拟机、虚拟硬盘、虚拟交换机、虚拟网卡等核心资源的自动化管理。

- 支持主要功能模块：虚拟机生命周期管理、虚拟硬盘操作、虚拟交换机与网络适配器管理等
- 适用场景：自动化运维、云平台集成、批量虚拟化资源管理、DevOps 工具链扩展等
- 依赖环境：
  - 仅支持 Windows 平台（需开启 Hyper-V 角色）
  - Go 1.18 及以上版本

## 安装

```bash
go get github.com/rokukoo/hyperv
```

> 仅支持 Windows，需提前在宿主机启用 Hyper-V。

## 功能

### 物理网卡 NetworkAdapter

物理网卡（Physical Network Adapter）是计算机硬件中用于实现网络连接的关键组件，通常以网卡（NIC, Network Interface Card）的形式存在于主机或服务器中。

在 Hyper-V 虚拟化环境中，物理网卡不仅负责主机本身的网络通信，还常常作为虚拟交换机（Virtual Switch）的底层承载，实现虚拟机与外部网络之间的数据转发。

常用方法:

```go
// ListAvailablePhysicalNetworkAdapters 列出所有可用的物理网络适配器
func ListAvailablePhysicalNetworkAdapters() ([]string, error)

// FindNetworkAdapterByName 根据名称查询网络适配器
// Not Implemented !
func FindNetworkAdapterByName(name string)

// EnableNetworkAdapter 启用网络适配器
// Not Implemented !
func EnableNetworkAdapter(name string) error

// DisableNetworkAdapter 禁用网络适配器
// Not Implemented !
func DisableNetworkAdapter(name string) error

// ConfigureNetworkAdapter 配置网络适配器
// Not Implemented !
func ConfigureNetworkAdapter(
  ipAddress string[],
  subnetMask string[],
  defaultGateway string[],
  dnsServers string[]
) error
```

### 虚拟机 VirtualMachine

虚拟机（Virtual Machine，简称 VM）是一种通过软件模拟的计算机系统，能够在物理主机上运行多个相互隔离的操作系统实例。每台虚拟机都拥有独立的 CPU、内存、存储和网络资源，用户可以像操作真实物理服务器一样对其进行管理和使用。

常用方法:

```go
// CreateVirtualMachine 创建虚拟机
func CreateVirtualMachine(name string, savePath string, cpuCoreCount int, memorySize int) (*VirtualMachine, error)

// DestroyVirtualMachineByName 根据名称销毁虚拟机
func DestroyVirtualMachineByName(name string, del bool) (ok bool, err error)

// DeleteVirtualMachineByName 根据名称删除虚拟机
func DeleteVirtualMachineByName(name string) (ok bool, err error)

// Start 启动虚拟机
func (vm *VirtualMachine) Start() error

// Stop 停止虚拟机
func (vm *VirtualMachine) Stop(force bool) error

// Shutdown 正常关闭虚拟机
func (vm *VirtualMachine) Shutdown() error

// ForceStop 强制停止虚拟机
func (vm *VirtualMachine) ForceStop() error

// Reboot 重启虚拟机
func (vm *VirtualMachine) Reboot(force bool) error

// ForceReboot 强制重启虚拟机
func (vm *VirtualMachine) ForceReboot() error

// Suspend 挂起虚拟机
// 由于 Hyper-V 平台原生挂起功能并非真正意义上的挂起, 而是保存虚拟机的状态, 因此这里的挂起操作实际上是保存虚拟机的状态
// 保存虚拟机的状态后, 可以通过 Resume 恢复虚拟机的运行
func (vm *VirtualMachine) Suspend() error

// Resume 恢复虚拟机
func (vm *VirtualMachine) Resume() error

// Snapshot 快照虚拟机
func (vm *VirtualMachine) Snapshot() error

// ModifyVirtualMachineSpecByName 根据虚拟机名称修改虚拟机规格
func ModifyVirtualMachineSpecByName(name string, cpuCoreCount int, memorySize int) (ok bool, err error)

// ModifyInternalIPv4Address 根据虚拟机名称修改IP地址
// Not Implemented !
func ModifyInternalIPv4Address() (ok bool, err error)

// FindVirtualMachineByName 根据虚拟机名称获取虚拟机
func FindVirtualMachineByName(vmName string) ([]*VirtualMachine, error)

// FirstVirtualMachineByName 根据虚拟机名称获取第一个虚拟机
func FirstVirtualMachineByName(vmName string) (*VirtualMachine, error)

// GetKvpItem 获取键值对
// Not Implemented !
func (vm *VirtualMachine) GetKvpItem()

// SetKvpItem 设置键值对
// Not Implemented !
func (vm *VirtualMachine) SetKvpItem()

// ListKvpItems 获取所有键值对
// Not Implemented !
func (vm *VirtualMachine) ListKvpItems()
```

### VirtualHardDisk

虚拟硬盘（Virtual Hard Disk，简称 VHD）是一种以文件形式存在的虚拟化存储设备，能够模拟真实物理硬盘的功能。虚拟硬盘广泛应用于虚拟机环境中，为虚拟机提供独立的存储空间，实现操作系统、应用程序和数据的隔离与管理。

在 Hyper-V 虚拟化平台中，虚拟硬盘支持多种类型（如系统盘、数据盘），可灵活挂载到不同的虚拟机上。通过 hypervctl，用户可以自动化完成虚拟硬盘的创建、删除、挂载、卸载、扩容等操作，并支持获取虚拟硬盘的详细信息（如名称、类型、容量、使用情况、路径等）。

```go
// CreateVirtualHardDisk 创建虚拟硬盘
func CreateVirtualHardDisk(path string, sizeGiB float64) (vhd *VirtualHardDisk, err error)

// DeleteVirtualHardDiskByPath 根据路径删除虚拟硬盘
func DeleteVirtualHardDiskByPath(path string) (ok bool, err error)

// Resize 调整虚拟硬盘大小
func (vhd *VirtualHardDisk) Resize(newSizeGiB float64) (ok bool, err error)

// AttachToByName 根据虚拟机名称挂载虚拟硬盘
func (vhd *VirtualHardDisk) AttachToByName(vmName string) (ok bool, err error)

// GetVirtualHardDiskByPath 根据路径获取虚拟硬盘信息
func GetVirtualHardDiskByPath(path string) (*VirtualHardDisk, error)
```

### VirtualSwitch

虚拟交换机（Virtual Switch）是 Hyper-V 中的一个重要网络组件，用于为虚拟机提供网络连接功能。它可以将多个虚拟网络适配器连接在一起，并根据不同的类型提供不同的网络连接方式。

Hyper-V 支持四种类型的虚拟交换机:

- External(外部): 可以让虚拟机通过物理网卡访问外部网络
- Internal(内部): 可以让虚拟机与宿主机及其他虚拟机进行通信
- Private(私有): 只允许虚拟机之间进行通信
- Bridge(桥接): 可以让虚拟机直接访问物理网络,类似于 External 类型

通过 hypervctl，用户可以创建、删除和管理不同类型的虚拟交换机，并可以修改虚拟交换机的类型。同时还支持查询虚拟交换机的详细信息，如名称、类型等。

```go
// CreateVirtualSwitch 创建虚拟交换机
// switchType: "External" | "Internal" | "Private" | "Bridge"
// physicalAdapterName 仅在 External/Bridge 类型下需要
func CreateVirtualSwitch(name string, switchType string, physicalAdapterName string) (*VirtualSwitch, error)

// DeleteVirtualSwitchByName 根据名称删除虚拟交换机
func DeleteVirtualSwitchByName(name string) (ok bool, err error)

// ChangeVirtualSwitchTypeByName 根据名称修改虚拟交换机类型
func ChangeVirtualSwitchTypeByName(name string, switchType VirtualSwitchType, adapter *string) error

// FirstVirtualSwitchByName 根据名称获取第一个虚拟交换机
func FirstVirtualSwitchByName(name string) (*VirtualSwitch, error)

// GetVirtualSwitchTypeByName 根据名称获取虚拟交换机类型
func GetVirtualSwitchTypeByName(name string) (VirtualSwitchType, error)

// ListVirtualSwitches 列出所有虚拟交换机
// Not Implemented !
func ListVirtualSwitches() ([]*VirtualSwitch, error)
```

### VirtualNetworkAdapter

虚拟网络适配器（Virtual Network Adapter）是虚拟机中的网络接口设备,用于为虚拟机提供网络连接功能。每个虚拟机可以配置多个虚拟网络适配器,并可以连接到不同的虚拟交换机上。

```go
// AddVirtualNetworkAdapter 添加虚拟网络适配器
func (vm *VirtualMachine) AddVirtualNetworkAdapter(vna *VirtualNetworkAdapter) (err error)

// DeleteVirtualNetworkAdapterByName 根据名称删除虚拟网卡
func (vm *VirtualMachine) RemoveVirtualNetworkAdapter(name string) (err error)

// SetBandwidth 设置虚拟网络适配器的带宽
func (vna *VirtualNetworkAdapter) SetBandwidth(limitBandwidthMbps, reserveBandwidthMbps float64) (err error)

// DisableBandwidthLimit 禁用虚拟网络适配器的带宽限制
func (vna *VirtualNetworkAdapter) DisableBandwidthLimit() (err error)

// SetMacAddress 设置虚拟网络适配器自动申请物理地址
// Not Implemented !
func (vna *VirtualNetworkAdapter) AutoMacAddress() (macAddress string, err error)

// SetMacAddress 设置虚拟网络适配器的物理地址
// Not Implemented !
func (vna *VirtualNetworkAdapter) SetMacAddress(macAddress string) (err error)

// GetMacAddress 获取虚拟网络适配器的物理地址
// Not Implemented !
func (vna *VirtualNetworkAdapter) GetMacAddress() (err error)

// EnableVirtualNetworkAdapterVlan 启用虚拟网卡VLAN并设置VLAN ID
// Not Implemented !
func (vna *VirtualNetworkAdapter) EnableVlan(adapterName string, vlanId int) (ok bool, err error)

// DisableVirtualNetworkAdapterVlan 禁用虚拟网卡VLAN
// Not Implemented !
func (vna *VirtualNetworkAdapter) DisableVlan(adapterName string) (ok bool, err error)

// ConnectByName 连接虚拟网络适配器到虚拟交换机
func (vna *VirtualNetworkAdapter) ConnectByName(vswName string) (bool, error)

// DisConnect 断开虚拟网络适配器与虚拟交换机的连接
func (vna *VirtualNetworkAdapter) DisConnect() (err error)

// ModifyConfiguration 修改虚拟网络适配器的配置
func (vna *VirtualNetworkAdapter) ModifyConfiguration(
	ipV4Address, subnetMask, defaultGateway, dnsServer []string,
) (err error)

// FindVirtualNetworkAdapterByName 根据名称查找虚拟网络适配器
func FindVirtualNetworkAdapterByName(name string) (virtualNetworkAdapters []*VirtualNetworkAdapter, err error)

// FirstVirtualNetworkAdapterByName 根据名称查找第一个虚拟网络适配器
func FirstVirtualNetworkAdapterByName(name string) (virtualNetworkAdapter *VirtualNetworkAdapter, err error)
```

### Cluster

// Not Implemented !
