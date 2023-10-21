<template>
  <div style="margin-top: -15px;text-align: left;user-select: text">
    <h3 style="color: black">使用说明</h3>
    <div style="font-size: 13px;">
      <div style="margin-bottom: 5px">1、单独间隔时间只有连发模式下生效</div>
      <div style="margin-bottom: 5px">2、设置单独间隔后，要取消需要删除重加</div>
      <div style="margin-bottom: 5px">3、按压模式尚在测试，有问题可以反馈到<el-link @click="clickWeb" style="font-size: 12px;vertical-align: text-bottom">评论区</el-link></div>
      <div >4、开启停止同时为滚轮上或下，有秒开秒关问题，悠着点</div>
    </div>
    <h3 style="color: black">更新内容</h3>
    <div style="font-size: 13px;">
      <div style="margin-bottom: 5px">1、开启、停止键位可设置为同一个</div>
      <div style="margin-bottom: 5px">2、暂停按键可以设置为空，在最下面</div>
      <div style="margin-bottom: 5px">3、声音提示关闭可保存配置</div>
      <div style="margin-bottom: 5px">4、键之间可以自行设定不同延时</div>
      <div style="margin-bottom: 5px">5、添加鼠标中键</div>
      <div style="margin-bottom: 5px">6、新增最小化到任务栏</div>
    </div>
    <div style="color: green;font-size: 12px;position: absolute;bottom: 10px;">当前版本：2023.10.21</div>
  </div>
</template>
<script lang="ts" setup>
  import {ElMessage} from "element-plus";

  const clickWeb = ()=>{
    navigator.clipboard.writeText("https://www.jx3box.com/tool/67045")
        .then(function() {
          ElMessage.success({
            message: '地址复制成功',
            type: 'success',
            offset: 50
          })
        })
        .catch(function(err) {
          console.error("复制到剪贴板失败: ", err);
        });
  }
</script>