<template>
  <el-button class="input-btn" type="primary" style="width: 100%" :disabled="isRunning" @click="settings"><el-icon><Tools /></el-icon>更多设置</el-button>
  <el-drawer
      v-model="drawer"
      title="设置"
      direction="rtl"
      size="230"
  >
    <div style="color: black;margin-top: -20px">
      <div>
        <div style="color: #409EFF;font-size: 13px"><el-icon style="vertical-align: text-bottom"><BellFilled /></el-icon>声音音量</div>
        <el-slider v-model="volume" @change="changeVolume" />
      </div>
      <div style="margin-top: 10px">
        <div style="color: #409EFF;font-size: 13px;margin-bottom: 10px"><el-icon style="vertical-align: text-bottom"><VideoPlay /></el-icon>开启音效</div>
        <el-input
            v-model="startMp3"
            placeholder="默认开启音效"
            class="input-with-select"
            readonly
            style="width: 100%"
        >
          <template #append>
            <el-button @click="setMp3('start')">自定义</el-button>
          </template>
        </el-input>
      </div>
      <div style="margin-top: 10px">
        <div style="color: #409EFF;font-size: 13px;margin-bottom: 10px"><el-icon style="vertical-align: text-bottom"><VideoPause /></el-icon>关闭音效</div>
          <el-input
              v-model="stopMp3"
              placeholder="默认停止音效"
              readonly
              class="input-with-select"
              style="width: 100%"
          >
            <template #append>
              <el-button @click="setMp3('stop')">自定义</el-button>
            </template>
          </el-input>
        </div>
      <div style="margin-top: 10px">
        <div style="color: #409EFF;font-size: 13px"><el-icon style="vertical-align: text-bottom"><SwitchButton /></el-icon>关闭方式</div>
        <el-radio-group v-model="close" @change="changeClose" class="ml-4">
          <el-radio label="0" size="large">退出</el-radio>
          <el-radio label="1" size="large">隐藏</el-radio>
        </el-radio-group>
      </div>
      <div style="margin-top: 10px">
        <div style="color: #409EFF;font-size: 13px"><el-icon style="vertical-align: text-bottom"><CoffeeCup /></el-icon>暂停方式</div>
        <el-radio-group v-model="parse" @change="changeParse" class="ml-4">
          <el-radio label="0" size="large">点按</el-radio>
          <el-radio label="1" size="large">按压</el-radio>
        </el-radio-group>
      </div>
      <div style="margin-top: 10px">
        <div style="color: #409EFF;font-size: 13px"><el-icon style="vertical-align: text-bottom"><Timer /></el-icon>独立延迟</div>
        <el-radio-group v-model="disabled" @change="changeDisabled" class="ml-4">
          <el-radio label="0" size="large">关闭</el-radio>
          <el-radio label="1" size="large">开启</el-radio>
        </el-radio-group>
      </div>
      <div style="margin-top: 10px">
        <div style="color: #409EFF;font-size: 13px;margin-bottom: 3px"><el-icon style="vertical-align: text-bottom"><StarFilled /></el-icon>方案选择--<span style="cursor: pointer" @click="savePlan">点击保存方案</span></div>
        <div style="color: orange"><span style="cursor: pointer;font-size: 12px;margin-right: 70px" @click="importPlan">导入</span><span style="cursor: pointer;font-size: 12px" @click="exportPlan">导出</span></div>
        <el-table :data="plans" border  style="width: 100%;font-size: 12px;" height="121" empty-text="暂无方案">
          <el-table-column prop="name" label="方案" width="100" align="center" />
          <el-table-column prop="select" label="操作" align="center">
            <template #header>
              <el-link v-if="!isDel" type="warning" style="font-size: 12px" @click="isDel = true">删除</el-link>
              <el-link v-else type="warning" style="font-size: 12px" @click="isDel = false">取消</el-link>
            </template>
            <template #default="scope">
              <el-link v-if="!isDel" type="primary" style="font-size: 12px" @click="selectPlan(scope.row)">选择</el-link>
              <el-link v-else type="danger" style="font-size: 12px" @click="delPlan(scope.row)">删除</el-link>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </el-drawer>
</template>
<script lang="ts" setup>
  import {ref} from "vue";
  import {SelectMusic} from "../../../wailsjs/go/service/Music";
  import {
    BellFilled,
    CoffeeCup,
    StarFilled,
    SwitchButton,
    Timer,
    Tools,
    VideoPause,
    VideoPlay
  } from "@element-plus/icons-vue";
  import {ElMessage, ElMessageBox} from "element-plus";
  import {WindowReloadApp} from "../../../wailsjs/runtime";
  import {ExportPlans, ImportPlans} from "../../../wailsjs/go/service/Keyboard";

  defineProps({
    isRunning:Boolean,
  })

  const emit = defineEmits(['emitVolume','emitMusic','emitParseType','emitDisabled'])

  const drawer = ref(false)

  const plans = ref<any>([])
  const savePlan = ()=>{
    ElMessageBox.prompt('给这个方案起个名吧！', '保存方案', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern:
          /^[a-zA-Z0-9\u4e00-\u9fa5]{1,5}$/,
      inputErrorMessage: '最多5个字且只能中文、英文和数字',
    })
        .then(({ value }) => {
          for (let i = 0; i < plans.value.length; i++) {
            if(plans.value[i].name === value){
              ElMessage.warning("这个方案已经有了，换个名字吧！")
              return
            }
          }
          selectId.value = plans.value.length-1
          let start = localStorage.getItem("start")||""
          let stop = localStorage.getItem("stop")||""
          let parse = localStorage.getItem("parse")||""
          let keys = localStorage.getItem("keys")||""
          let model = localStorage.getItem("model")||""
          let ms = localStorage.getItem("ms")||""
          let parseType = localStorage.getItem("parseType")||""
          let disabled = localStorage.getItem("disabled")||""
          plans.value.push({name:value,start:start,stop:stop,parse:parse,keys:keys,model:model,ms:ms,parseType:parseType,disabled:disabled})
          localStorage.setItem("plans",JSON.stringify(plans.value))
        })
        .catch(() => {})
  }

  const exportPlan = () =>{
    ExportPlans(JSON.stringify(plans.value)).then(res=>{
      ElMessage.success(res)
    })
  }

  const importPlan = () =>{
    ElMessageBox.confirm(
        '导入操作将覆盖所有方案，确定？',
        '导入',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
        .then(() => {
            ImportPlans().then(res=>{
              if(res == "取消导入"){
                ElMessage.success(res)
              } else {
                plans.value = JSON.parse(res)
                localStorage.setItem("plans",JSON.stringify(plans.value))
                ElMessage.success("导入成功")
              }
            })
        })
        .catch(() => {})
  }

  const selectPlan = (val: any) =>{
    localStorage.setItem("start",val.start)
    localStorage.setItem("stop",val.stop)
    localStorage.setItem("parse",val.parse)
    localStorage.setItem("keys",val.keys)
    localStorage.setItem("model",val.model)
    localStorage.setItem("ms",val.ms)
    localStorage.setItem("parseType",val.parseType)
    localStorage.setItem("disabled",val.disabled)
    WindowReloadApp()
  }

  const delPlan = (val: any) => {
    for (let i = 0; i < plans.value.length; i++) {
      if(plans.value[i].name == val.name){
        plans.value.splice(i,1)
        localStorage.setItem("plans",JSON.stringify(plans.value))
      }
    }
  }

  const isDel = ref(false)

  const volume = ref()
  const startMp3 = ref()
  const stopMp3 = ref()
  const close = ref()
  const parse = ref()
  const disabled = ref()
  const selectId = ref(-1)

  const settings = ()=>{
    // localStorage.removeItem("plans")
    drawer.value = true
    volume.value = Number(localStorage.getItem('volume') || "50")
    startMp3.value = localStorage.getItem('startMusic') || ""
    stopMp3.value = localStorage.getItem('stopMusic') || ""
    close.value = localStorage.getItem('close') || ""
    parse.value = localStorage.getItem('parseType') || "0"
    disabled.value = localStorage.getItem('disabled') || "0"
    plans.value = JSON.parse(localStorage.getItem("plans")||"[]")
  }
  const changeVolume = ()=>{
    emit('emitVolume',volume.value)
  }
  const changeClose = ()=>{
    localStorage.setItem("close",close.value)
  }
  const changeParse = ()=>{
    emit('emitParseType',parse.value)
  }

  const changeDisabled = ()=>{
    emit('emitDisabled',disabled.value)
  }

  const setMp3 = async (type: string) => {
    await SelectMusic().then(res => {
      if (type === 'start') {
        startMp3.value = res
      }
      if (type === 'stop') {
        stopMp3.value = res
      }
      emit('emitMusic',startMp3.value,stopMp3.value)
    })
  }
</script>