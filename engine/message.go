/**
 * @Author: cmpeax Tang
 * @Date: 2021/3/10 21:03
 */

package engine

type MessageID uint16

//
const (
	MESSAGE_SEND_TO_MONITOR        MessageID = iota // 发送信息到防静电设备
	MESSAGE_SEND_TO_FACE_GUARD                      // 发送信息到人脸门禁
	MESSAGE_SEND_TO_NORMAL_GUARD                    // 发送信息到普通门禁
	MESSAGE_SEND_TO_ALL_GUARD                       // 发送信息到全部门禁设备
	MESSAGE_RECV_FROM_MONITOR                       // 接收到防静电设备信息
	MESSAGE_RECV_FROM_FACE_GUARD                    // 接收到人脸门禁信息
	MESSAGE_RECV_FROM_NORMAL_GUARD                  // 接收到普通门禁信息
)

var MapMessageKey2Value map[string]MessageID = map[string]MessageID{
	"MESSAGE_SEND_TO_MONITOR":        MESSAGE_SEND_TO_MONITOR,
	"MESSAGE_SEND_TO_FACE_GUARD":     MESSAGE_SEND_TO_FACE_GUARD,
	"MESSAGE_SEND_TO_NORMAL_GUARD":   MESSAGE_SEND_TO_NORMAL_GUARD,
	"MESSAGE_SEND_TO_ALL_GUARD":      MESSAGE_SEND_TO_ALL_GUARD,
	"MESSAGE_RECV_FROM_MONITOR":      MESSAGE_RECV_FROM_MONITOR,
	"MESSAGE_RECV_FROM_FACE_GUARD":   MESSAGE_RECV_FROM_FACE_GUARD,
	"MESSAGE_RECV_FROM_NORMAL_GUARD": MESSAGE_RECV_FROM_NORMAL_GUARD,
}
