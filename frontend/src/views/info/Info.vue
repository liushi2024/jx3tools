<template>
  <el-link type="warning" style="float:right;margin-top: -15px;margin-right: 10px;" @click="drawer=true">使用说明</el-link>
  <el-drawer v-model="drawer" :with-header="false" size="400" direction="btt" style="color: dimgray;">
    <div style="float: right">
      <el-button type="danger" @click="drawer = false">
        <el-icon class="el-icon--left"><CircleCloseFilled /></el-icon>
        <span>关闭</span>
      </el-button>
    </div>
    <div style="margin-top: -15px;text-align: left;user-select: text">
      <h3 style="color: black">使用说明</h3>
      <div style="font-size: 13px;">
        <div style="margin-bottom: 5px">1、有问题可以反馈到<el-link @click="clickWeb(0)" style="font-size: 12px;vertical-align: text-bottom">评论区</el-link></div>
        <div style="margin-bottom: 5px">2、单独延迟需在更多设置里开启独立延迟才生效</div>
        <div style="margin-bottom: 5px">3、角色多的可以用<el-link @click="clickWeb(1)" style="font-size: 12px;vertical-align: text-bottom">角色同步小工具</el-link></div>
      </div>
      <h3 style="color: black">注意事项</h3>
      <div style="font-size: 13px;">
        <div style="margin-bottom: 5px">1、多段模式不写了，有bug，懒</div>
      </div>
      <div style="color: green;font-size: 12px;position: absolute;bottom: 10px;">当前版本：2023.11.05</div>
    </div>
  </el-drawer>
</template>
<script lang="ts" setup>
  import {CircleCloseFilled} from "@element-plus/icons-vue";
  import {ref} from "vue";
  import {BrowserOpenURL} from "../../../wailsjs/runtime";


  const drawer = ref(false)

  const clickWeb = (val: number)=> {
    switch (val) {
      case 0: BrowserOpenURL("https://www.jx3box.com/tool/67045");break;
      case 1: BrowserOpenURL("https://www.jx3box.com/tool/67154");break;
      default:
    }

  }
</script>