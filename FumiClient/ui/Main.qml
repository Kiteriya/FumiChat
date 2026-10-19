import QtQuick
import QtQuick.Controls

ApplicationWindow {
    minimumHeight: 600
    minimumWidth: 800
    visible: true
    title: "FumiChat"
    SplitView{
        anchors.fill:parent
        orientation: Qt.Horizontal
        Rectangle{
            id: leftSideBar
            color: "#FFFFFF"
            SplitView.preferredWidth: 320
            SplitView.minimumWidth: 260
            SplitView.maximumWidth: 450
            ListView{
                id: chatList
                orientation: ListView.Vertical
                spacing: 5

            }
        }
    }
}