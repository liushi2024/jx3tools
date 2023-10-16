<template>
  <div class="header" :style="{backgroundColor:isRunning?'#f5222d':'#409EFF'}">
    <div style="margin-left: 10px;">
      <el-icon size="21" style="margin-top: 13px;margin-right: 15px;"><SwitchFilled /></el-icon>
      <div class="name">剑三小助手   -   <span style="font-size: 8px;">{{isRunning?"运行中":"未运行"}}</span></div>

    </div>
    <div v-if="true">
      <el-icon class="btn-sys" size="small" style="margin-right: 10px" @click="sys_event('min')"><Minus /></el-icon>
      <el-icon class="btn-sys" size="small" style="margin-right: 10px" @click="sys_event('close')"><Close /></el-icon>
    </div>
  </div>
  <div style="padding: 10px;overflow: hidden;display: flex;justify-content: space-between ">
    <div class="left-box">
      <div style="width: 100%;height: 50px;">
        <el-button class="input-btn" @click="preKeys" :disabled="isRunning">{{preKey==""?"点击输入按键":preKey}}</el-button>
      </div>
      <el-table
          class="table-box"
          :data="keyList"
          style="width: 100%;font-size: 12px"
          height="calc(100vh - 160px)"
      >
        <el-table-column width="55" label="启用" align="center">
          <template #default="scope">
            <el-switch
                @change="onSaveKey()"
                v-model="scope.row.used"
                inline-prompt
                :disabled="isRunning"
                active-text="是"
                inactive-text="否"
            />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="键符" width="50" align="center" show-overflow-tooltip/>
        <el-table-column prop="status" label="操作" align="center">
          <template #default="scope">
            <el-link type="danger" style="font-size: 13px" :disabled="isRunning" @click="onDelKey(scope.row)">删除</el-link>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="right-box">
      <el-button class="input-btn" style="margin-bottom: 35px;width: 100%" :disabled="isRunning" @click="addKey">添加</el-button>
      <el-divider content-position="left" style="width: 100%;margin: 0 auto 20px;">开启</el-divider>
      <el-select v-model="start" class="m-2" filterable placeholder="开启键" :disabled="isRunning" :key="1" size="large" style="width: 100%" @visible-change="filterKey" @change="onSaveHotKey('start')">
        <el-option
            v-for="item in keyOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
        />
      </el-select>
      <el-divider content-position="left" style="width: 100%;margin: 20px auto 20px;">停止</el-divider>
      <el-select v-model="stop" class="m-2" filterable placeholder="停止键" :disabled="isRunning" size="large" style="width: 100%;" @visible-change="filterKey" @change="onSaveHotKey('stop')">
        <el-option
            v-for="item in keyOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
        />
      </el-select>
      <el-divider content-position="left" style="width: 100%;margin: 20px auto 20px;">暂停</el-divider>
      <el-select v-model="parse" class="m-2" filterable placeholder="暂停键" :disabled="isRunning" size="large" style="width: 100%;" @visible-change="filterKey" @change="onSaveHotKey('parse')">
        <el-option
            v-for="item in keyOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
        />
      </el-select>
      <el-divider content-position="left" style="width: 100%;margin: 20px auto 20px;">模式</el-divider>
      <el-select v-model="model" class="m-2" filterable placeholder="模式" :disabled="isRunning" size="large" style="width: 100%;" @change="onSaveHotKey('model')">
        <el-option
            v-for="item in modelOption"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <el-divider content-position="left" style="width: 100%;margin: 20px auto 20px;" >间隔</el-divider>
      <el-input v-model="ms" placeholder="" :disabled="isRunning" type="number" style="width: 100%;" @change="onSaveHotKey('ms')">
        <template #suffix >
          <span>ms</span>
        </template>
      </el-input>
      <el-form-item label="声音提示" size="small" class="voice">
        <el-switch v-model="voice" :disabled="isRunning" @change="onSaveHotKey('voice')" style="margin-left: 20px"/>
      </el-form-item>
    </div>
  </div>
  <el-text class="mx-1" type="warning" style="float: left;margin-top: -15px;margin-left: 10px"><el-icon><InfoFilled /></el-icon>若无法使用，请<el-link @click="clear" type="danger" style="margin-top: -2.6px;margin-left: 2px">重置按键设置</el-link></el-text>
</template>
<script setup lang="ts">
import {onMounted, ref} from "vue";
import {ElLoading, ElMessage, ElMessageBox} from "element-plus";
import {Close, InfoFilled, Minus, SwitchFilled} from "@element-plus/icons-vue";
import {
  DllImport,
  SyncFrontKey,
  SyncFrontModel,
  SyncFrontMs,
  SyncFrontVoice,
} from "../../../wailsjs/go/service/Keyboard";
import {LogDebug, EventsOn, Quit,WindowMinimise} from "../../../wailsjs/runtime";
import {SyncFrontParse, SyncFrontStart, SyncFrontStop} from "../../../wailsjs/go/service/HotKey";

const vk:any = {'1': 201,'2': 202,'3': 203,'4': 204,'5': 205,'6': 206,'7': 207,'8': 208,'9': 209,'0': 210,
    'a': 401,'b': 505,'c': 503,'d': 403,'e': 303,'f': 404, 'g': 405,'h': 406,'i': 308,'j': 407, 'k': 408,
    'l': 409,'m': 507,'n': 506,'o': 309,'p': 310,'q': 301,'r': 304, 's': 402,'t': 305,'u': 307,'v': 504,
    'w': 302, 'x': 502,'y': 306,  'z': 501,   '[': 311, '.': 509, ',': 508, ']': 312, '/': 510,
    "'": 411, '\\': 313,  '`': 200,  '-': 211, '=': 212,  ';': 410,
    'F1':101,'F2':102,'F3':103,'F4':104,'F5':105,'F6':106,'F7':107,'F8':108,'F9':109,'F10':110,'F11':111,
    'CapsLock':400,'Shift':500,'Ctrl':600,'Alt':602,'Tab':300,"ScrollLock":701,'pause':702,'Insert':703,'Home':704,'PageUp':705,'Delete':706,'End':707,'PageDown':708,
    'ArrowUp':709,'ArrowLeft':710,'ArrowDown':711,'ArrowRight':712}
const vq:any = {'0':'48',
  '1':'49',
  '2':'50',
  '3':'51',
  '4':'52',
  '5':'53',
  '6':'54',
  '7':'55',
  '8':'56',
  '9':'57',
  'A':'65',
  'B':'66',
  'C':'67',
  'D':'68',
  'E':'69',
  'F':'70',
  'G':'71',
  'H':'72',
  'I':'73',
  'J':'74',
  'K':'75',
  'L':'76',
  'M':'77',
  'N':'78',
  'O':'79',
  'P':'80',
  'Q':'81',
  'R':'82',
  'S':'83',
  'T':'84',
  'U':'85',
  'V':'86',
  'W':'87',
  'X':'88',
  'Y':'89',
  'Z':'90',
  'F1':'112',
  'F2':'113',
  'F3':'114',
  'F4':'115',
  'F5':'116',
  'F6':'117',
  'F7':'118',
  'F8':'119',
  'F9':'120',
  'F10':'121',
  'F11':'122',
  'Left':'37',
  'Right':'38',
  'Up':'39',
  'Down':'40',
  'Alt':'1',
  'Ctrl':'2',
  'Shift':'4',
  '鼠标侧键1': '904',
  '鼠标侧键2': '905',
  '滚轮上': '906',
  '滚轮下': '908'}
  const start = ref("")
  const stop = ref("")
  const parse = ref("")
  const model = ref(0)
  const keyList = ref<any[]>([])
  const ms = ref(0)
  const keyOptions = ref<any[]>([])
  const modelOption = ref([{label:'顺序模式',value:0},{label: '连发模式',value:1}])
  const preKey = ref("")
  const voice = ref(false)
  const isRunning = ref(false)

  const sys_event = (val: any)=>{
    switch(val){
      case 'min':WindowMinimise();break;
      case 'close':Quit();Quit();break;
      default:console.log("click btn_sys error");break;
    }
  }
  const preKeys = ()=> {
    const loading = ElLoading.service({
      lock: true,
      text: '等待输入按键，按esc退出',
      background: 'rgba(255, 255, 255, 0.9)',
    })
    document.onkeydown = function (event) {
      document.onkeydown = null
      loading.close()
      if(event.key == "Escape"){
        return;
      }
      preKey.value = event.key
    }
  }

  const addKey = ()=>{
    if(vk[preKey.value]==undefined){
      ElMessage.warning(preKey.value+"键暂不支持！")
      return;
    }
    for (let i = 0; i < keyList.value.length; i++) {
      if(keyList.value[i].name == preKey.value){
        ElMessage.warning(preKey.value+"键冲突！")
        return ;
      }
    }
    if(vk[preKey.value] == start.value){
      ElMessage.warning("此键已被开启键占用！")
      return;
    }
    if(vk[preKey.value] == stop.value){
      ElMessage.warning("此键已被停止键占用！")
      return;
    }
    keyList.value.push({name:preKey.value,value:vk[preKey.value],used:true})
    onSaveKey()
  }

  const  filterKey=(event: any)=>{
    if(event){
      keyOptions.value = []
      for (let key in vq) {
        let disabled = false;
        for (let i = 0; i < keyList.value.length; i++) {
          if(keyList.value[i].name == key){
            disabled = true;
          }
        }
        if(start.value == vq[key]||stop.value == vq[key]||parse.value == vq[key]){
          disabled = true
        }
        keyOptions.value.push({label:key,value:vq[key],disabled:disabled})
      }
    }
  }

  const onSaveHotKey = async (val: string) => {
    let value: any = ""
    if (val == "start") {
      value = start.value
      await SyncFrontStart(Number(value))
    } else if (val == "stop") {
      value = stop.value
      await SyncFrontStop(Number(value))
    } else if (val == "parse") {
      value = parse.value
      await SyncFrontParse(Number(value))
    } else if (val == "ms") {
      if (ms.value <= 50||ms.value == null){
        ElMessage.warning("间隔最小为50")
        ms.value = 50
      }
      value = ms.value
      await SyncFrontMs(Number(value))
    } else if (val == "voice") {
      value = voice.value
      await SyncFrontVoice(Boolean(value))
    } else if (val == "model") {
      value = model.value
      await SyncFrontModel(Number(value))
    } else {
      return;
    }
    localStorage.setItem(val, value)
  }

  const onSaveKey=async () => {
    await initData()
    localStorage.setItem("keys",JSON.stringify(keyList.value))
  }
  const onDelKey=async (value: { name: any; }) => {
    for (let i = 0; i < keyList.value.length; i++) {
      if (keyList.value[i].name == value.name) {
        keyList.value.splice(i, 1)
      }
    }
    await initData()
    localStorage.setItem("keys", JSON.stringify(keyList.value))
  }

  const initParam = (key: string, value: string) =>{
    let keyStr =  localStorage.getItem(key)
    if(keyStr == undefined){
      keyStr = value
      localStorage.setItem(key,keyStr)
    }
    return keyStr
  }

  const initData = async () => {
    let temp = []
    for (let i = 0; i < keyList.value.length; i++) {
      if (Boolean(keyList.value[i].used)) {
        temp.push(keyList.value[i].value)
      }
    }
    await SyncFrontKey(temp)
  }

  const clear=()=>{
    ElMessageBox.confirm(
        '将清除所有按键设置，重置后请重新打开程序',
        '重置',
        {
          confirmButtonText: '重置',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
        .then(() => {
          localStorage.clear();
          Quit()
        })
        .catch(() => {})
  }

  onMounted(async ()=>{
    let res = await DllImport()
    res == "success"?LogDebug("驱动注入成功"):ElMessage.warning(res);
    const loading = ElLoading.service({
      lock: true,
      text: '正在加载，请稍后',
      background: 'rgba(255, 255, 255, 0.9)',
    })
    start.value = initParam("start",vq['F9'])
    stop.value = initParam("stop",vq['F10'])
    parse.value = initParam("parse",vq['F11'])
    model.value = Number(initParam("model","0"))
    ms.value = Number(initParam("ms","50"))
    voice.value = Boolean(initParam("voice","true"))
    let keys =  localStorage.getItem("keys")
    if(keys == undefined){
      keyList.value = []
    } else {
      keyList.value = JSON.parse(keys)
    }
    filterKey(true)
    await SyncFrontStart(Number(start.value))
    await SyncFrontStop(Number(stop.value))
    await SyncFrontParse(Number(parse.value))
    await SyncFrontMs(ms.value)
    await SyncFrontModel(model.value)
    await SyncFrontVoice(voice.value)
    await initData()
    EventsOn("start-thread",()=>{
      isRunning.value = true
    })
    EventsOn("stop-thread",()=>{
      isRunning.value = false
    })
    loading.close()
  })
</script>
<style lang="less" scoped>
  @import "../../assets/style/theme";
  .header{
    --wails-draggable:drag;
    width: 100%;
    height: 48px;
    color: white;
    line-height: 48px;
    display: flex;
    justify-content: space-between;
  }
  .name{
    display: inline-block;
    font-size: 14px;
    user-select: none;
    vertical-align: top;
    margin-left: -10px;
  }
  .btn-sys{
    display: inline-block;
    color:white;
    :hover{
      cursor: pointer;
      color: black;
    }
  }
  .left-box{
    //display: inline-block;
    width: 60%;
    .input-btn{
      width: 100%;
      box-shadow: @box-shadow-base;
    }
    .table-box{
      width: 100%;
      //height: calc(100vh - 28px);
      border-radius: 5px !important;
      overflow-y: hidden;
      box-shadow: @box-shadow-base;

    }
  }
  .right-box{
    //display: inline-block;
    width: 35%;
    //float: right;
    .input-btn{
      box-shadow: @box-shadow-base;
    }
    .voice{
      font-size: 16px;
      margin-top: 14px;
      color: @primary-color
    }
  }
</style>