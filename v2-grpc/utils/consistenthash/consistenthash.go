package consistenthash

// 一致性哈希算法的简单实现
// 节点包括服务器节点和数据节点。服务器节点包括物理节点和对应的虚拟节点，但是哈希环的ring只存储虚拟节点，并且用map存储虚拟节点和物理节点的映射关系
// 使用虚拟节点（为了让数据更均匀地分布在哈希环上），支持加权分配（即不同服务器节点可以有不同数量的虚拟节点）
// 包括加入节点并更新哈希环、删除节点并更新哈希环、根据key获取节点。但是，不包括数据迁移

import (
	"hash/fnv"
	"sort"
)

const (
	WEIGHT_DEFAULT = iota + 1// 默认权重
	WEIGHT_2
)

var WEIGHT_VIRTUAL_NODE_MAP = map[int]int {
	WEIGHT_DEFAULT: 100,
	WEIGHT_2: WEIGHT_DEFAULT * WEIGHT_2,
}

type Node struct {
	Name string // 节点名称
}

type PhysicalServerNode struct {
	Node
	Weight int // 节点权重
}

type VirtualNode struct {
	Node
}

type HashRing struct {
	length int64 // 哈希环的长度
	virtualNodeMap map[string]*PhysicalServerNode // 虚拟节点到物理节点的映射
	ring map[int64]VirtualNode // 哈希环，只存储虚拟节点，不存储物理节点
	sortedHashes []int64 // 哈希环上虚拟节点的哈希值，排序后用于查找
}

func (hr *HashRing) NewHashRing(length int64) *HashRing {
	return &HashRing{
		length: length,
		// 需要初始化吗？
	}
}

// 添加物理节点，根据权重生成对应数量的虚拟节点。只更新哈希环的拓扑结构，不执行实际的数据迁移
func (hr *HashRing) AddPhysicalServerNode(node *PhysicalServerNode){
	virtualNodeCount := WEIGHT_VIRTUAL_NODE_MAP[node.Weight]

	// 创建虚拟节点并加入
	for i := 0; i < virtualNodeCount; i++ {
		vnodeName := node.Name + "#" + string(i)
		vnodeHash := hashString(vnodeName) % hr.length

		vnode := VirtualNode{Node{Name: vnodeName}}

		hr.ring[vnodeHash] = vnode
		hr.virtualNodeMap[vnodeName] = node
		hr.sortedHashes = append(hr.sortedHashes, vnodeHash)
	}

	// 对哈希值进行排序
	sort.Slice(hr.sortedHashes, func(i, j int) bool {
		return hr.sortedHashes[i] < hr.sortedHashes[j]
	})
}

// 删除物理节点，并删除对应的虚拟节点
func (hr *HashRing) DeletePhysicalServerNode(node *PhysicalServerNode) {
	virtualNodeCount := WEIGHT_VIRTUAL_NODE_MAP[node.Weight]
	// 删除虚拟节点
	for i := 0; i < virtualNodeCount; i++ {
		vnodeName := node.Name + "#" + string(i)
		vnodeHash := hashString(vnodeName) % hr.length
		delete(hr.ring, vnodeHash)
		delete(hr.virtualNodeMap, vnodeName)
	}

	// 大量虚拟节点被删除，重新构造sortedHashes
	hr.sortedHashes = hr.sortedHashes[:0] // 通过截取的办法清空切片，但是保留底层数组和容量，避免多次扩容
	for h := range hr.ring {
		hr.sortedHashes = append(hr.sortedHashes, h)
	}
	sort.Slice(hr.sortedHashes, func(i, j int) bool {
		return hr.sortedHashes[i] < hr.sortedHashes[j]
	})
}

// 根据key获取物理节点，即给定数据应该在哪里物理节点上存储。函数名隐去虚拟节点的信息
func (hr *HashRing) GetNode(key string) *PhysicalServerNode {
	if len(hr.ring) == 0 {
		return nil
	}

	// 计算key的哈希值
	keyHash := hashString(key) % hr.length

	// 找到第一个大于等于keyHash的虚拟节点
	idx := sort.Search(len(hr.sortedHashes), func(i int) bool {
		return hr.sortedHashes[i] >= keyHash
	})

	// 如果没有找到，说明keyHash大于所有虚拟节点的哈希值，应该返回第一个虚拟节点（环形）
	if idx == len(hr.sortedHashes) {
		idx = 0
	}

	// 获取对应的物理节点
	physicalNode := hr.virtualNodeMap[hr.ring[hr.sortedHashes[idx]].Name]
	return physicalNode
}
// ---------------------- Utils ----------------------

// 计算字符串的哈希值
func hashString(key string) int64 {
	// 使用FNV-1a，计算哈希值（快速、简单的非加密哈希函数）
	h := fnv.New64a()
	h.Write([]byte(key))
	return int64(h.Sum64())
}