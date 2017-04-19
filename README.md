集群操作

Mirror 构建镜像 分Copy 和 Sync 
Move(dataset, source superset,target superset) 数据集在超级直接移动

定义。去SuperSet 是数据节点，作为replica 的最小单元

Dataset 数据存储的最小单元，里面都是kv集合

dataset 可以在superset 之间移动

不同的superset 之间建立mirror 关系。 建立mirror关系的set 保持一致性。由于技术限制，这里设计为最终一致

第一版本:
每个dataset 都有独立的binlog 通过这些binlog可以构建mirror

superset之间打通数据通道，不同的superset 只需要建立一条数据通道。superset 里面所有dataset 都通过这个通道进行同步


dataset 实现接口:mirror(exporter importer)     kv

superset    CopyAndSync(stream)  拷贝以及同步数据。 请求发送所有dataset 的名字以及binlog sequence 。 superset 通过这些信息，决定是copy 还是 sync 


一致性。一致性同步采用binlog 方式，从主拉取一次拷贝，再由binlog 同步更新。

状态有 本机状态  拷贝状态   镜像状态

复制为0 时，状态为本机状态
当开始复制，为复制状态(复制状态不可中断，中断需要重新拷贝



实现例子

一致性哈西实现:  每个dataset 为一个slot


基于slot 的merkle tree.用于对数据


