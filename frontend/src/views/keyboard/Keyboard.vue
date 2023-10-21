<template>
  <audio id="startMp3" controls style="display: none">
    <source src="../../../public/start.mp3" type="audio/mpeg">
  </audio>
  <audio id="stopMp3" controls style="display: none">
    <source src="../../../public/stop.mp3" type="audio/mpeg">
  </audio>
  <div class="header" :style="{backgroundColor:isRunning?'#f5222d':'#409EFF',}">
    <div style="margin-left: 10px;">
      <el-icon size="21" style="margin-top: 13px;margin-right: 15px;"><SwitchFilled /></el-icon>
      <div class="name">按键小助手   -   <span style="font-size: 8px;">{{isRunning?"运行中":"未运行"}}</span></div>
    </div>
    <div v-if="true">
      <el-icon class="btn-sys" size="small" style="margin-right: 10px" @click="sys_event('min')"><Minus /></el-icon>
      <el-icon class="btn-sys" size="small" style="margin-right: 10px" @click="sys_event('hide')"><Close /></el-icon>
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
        <el-table-column width="52" label="启用" align="center">
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
        <el-table-column prop="name" label="键符" width="50" align="center" show-overflow-tooltip>
          <template #default="scope">
            <el-link>{{scope.row.name}}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="操作" align="center"  show-overflow-tooltip>
          <template #header>
            <el-link type="warning" style="font-size: 12px" v-if="!isDel" @click="isDel = true">删除</el-link>
            <el-link type="warning" style="font-size: 12px;" v-else @click="isDel = false">取消</el-link>
          </template>
          <template #default="scope">
            <el-link type="danger" v-if="isDel" style="font-size: 12px" :disabled="isRunning" @click="onDelKey(scope.row)">删除</el-link>
            <el-link type="success" v-else style="font-size: 12px" :disabled="isRunning" @click="changeMs(scope.row)">{{scope.row.key_ms==0?'未配置':scope.row.key_ms+'ms'}} </el-link>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="right-box">
      <el-button class="input-btn" style="margin-bottom: 35px;width: 100%" :disabled="isRunning" @click="addKey">添加</el-button>
      <el-divider content-position="left" style="width: 100%;margin: 0 auto 20px;">开启</el-divider>
      <el-select placement="bottom" v-model="start" popper-class="custom-popper-class" class="m-2" filterable placeholder="开启键" :disabled="isRunning" :key="1" size="large" style="width: 100%" @visible-change="filterKey" @change="onSaveHotKey('start')">
        <el-option
            v-for="item in keyOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
        />
      </el-select>
      <div v-show="model != 2">
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
      </div>
      <div v-show="model == 2" style="height: 144px;font-size: 14px">
          <div style="margin-top: 20px;color: red">按压模式暂只支持键盘区！</div>
          <div style="margin-top: 20px;color: red">Ctrl、Alt、Shift可能存在冲突</div>
      </div>
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
  <el-link style="float:right;margin-top: -15px;margin-right: 10px" @click="drawer=true">使用说明</el-link>
  <el-drawer v-model="drawer" :with-header="false" size="400" direction="btt" style="color: dimgray;">
    <div style="float: right">
      <el-button type="danger" @click="drawer = false">
        <el-icon class="el-icon--left"><CircleCloseFilled /></el-icon>
        <span>关闭</span>
      </el-button>
    </div>
    <Info></Info>
  </el-drawer>
</template>
<script setup lang="ts">
import {onMounted, ref} from "vue";
import {ElLoading, ElMessage, ElMessageBox} from "element-plus";
import {CircleCloseFilled, Close, InfoFilled, Minus, SwitchFilled} from "@element-plus/icons-vue";
import {
  DllImport,
  SyncFrontKey,
  SyncFrontModel,
  SyncFrontMs,
  SyncFrontVoice,
} from "../../../wailsjs/go/service/Keyboard";
import {
  LogDebug,
  EventsOn,
  Quit,
  WindowHide,
  WindowMinimise,
  WindowSetTitle,
  WindowGetPosition, WindowSetPosition
} from "../../../wailsjs/runtime";
import {SyncFrontParse, SyncFrontStart, SyncFrontStop} from "../../../wailsjs/go/service/HotKey";
import Info from "../info/Info.vue";

const vk:any = {'1': 201,'2': 202,'3': 203,'4': 204,'5': 205,'6': 206,'7': 207,'8': 208,'9': 209,'0': 210,
    'a': 401,'b': 505,'c': 503,'d': 403,'e': 303,'f': 404, 'g': 405,'h': 406,'i': 308,'j': 407, 'k': 408,
    'l': 409,'m': 507,'n': 506,'o': 309,'p': 310,'q': 301,'r': 304, 's': 402,'t': 305,'u': 307,'v': 504,
    'w': 302, 'x': 502,'y': 306,  'z': 501,   '[': 311, '.': 509, ',': 508, ']': 312, '/': 510,
    "'": 411, '\\': 313,  '`': 200,  '-': 211, '=': 212,  ';': 410,
    'F1':101,'F2':102,'F3':103,'F4':104,'F5':105,'F6':106,'F7':107,'F8':108,'F9':109,'F10':110,'F11':111,
    'CapsLock':400,'Shift':500,'Ctrl':600,'Alt':602,'Tab':300,"ScrollLock":701,'pause':702,'Insert':703,'Home':704,'PageUp':705,'Delete':706,'End':707,'PageDown':708,
    'ArrowUp':709,'ArrowLeft':710,'ArrowDown':711,'ArrowRight':712}
const vq:any = {
  '0':'48', '1':'49', '2':'50', '3':'51', '4':'52', '5':'53', '6':'54', '7':'55', '8':'56', '9':'57',
  'A':'65', 'B':'66', 'C':'67', 'D':'68', 'E':'69', 'F':'70', 'G':'71', 'H':'72', 'I':'73', 'J':'74', 'K':'75', 'L':'76', 'M':'77', 'N':'78',
  'O':'79', 'P':'80', 'Q':'81', 'R':'82', 'S':'83', 'T':'84', 'U':'85', 'V':'86', 'W':'87', 'X':'88', 'Y':'89', 'Z':'90',
  'F1':'112', 'F2':'113', 'F3':'114', 'F4':'115', 'F5':'116', 'F6':'117', 'F7':'118', 'F8':'119', 'F9':'120', 'F10':'121', 'F11':'122',
  '-':'189', '+':'187', '~':'192', '\\':'220', ';':'186', '\'(引号)':'222', '[':'219', ']':'221', '，':'188', '。':'190', '/':'191',
  'Left':'37', 'Right':'39', 'Up':'38', 'Down':'40',
  'Left Alt':'164', 'Right Alt':'165', 'Left Ctrl':'162', 'Right Ctrl':'163', 'Left Shift':'160', 'Right Shift':'161',
  'PrtSc':'44', 'ScrollLock':'145', 'Pause':'19', 'Insert':'45', 'Delete':'46', 'Home':'36', 'End':'35', 'PgUp':'33', 'PgDown':'34',
  '鼠标侧键1': '904', '鼠标侧键2': '905', '滚轮上': '906', '滚轮下': '908', '空': "99999"}

  const isDel = ref(false)
  const drawer = ref(false)
  const start = ref("")
  const stop = ref("")
  const parse = ref("")
  const model = ref(0)
  const keyList = ref<any[]>([])
  const ms = ref(0)
  const keyOptions = ref<any[]>([])

  const modelOption = ref([
    {label: '顺序模式',value: 0},
    {label: '连发模式',value: 1},
    {label: '按压模式',value: 2}
  ])

  const preKey = ref("")
  const voice = ref(false)
  const isRunning = ref(false)

  const initWindowPosition = () => {
    let x = localStorage.getItem("x") || ""
    let y = localStorage.getItem("y") || ""
    if (!(x == "" || y == "")) {
      WindowSetPosition(Number(x), Number(y))
    }
  }
  initWindowPosition()
  /**
   * header点击事件
   * @param val
   */
  const sys_event = (val: any)=>{
    switch(val){
      case 'min':WindowMinimise();break;
      case 'hide':WindowHide();WindowGetPosition().then(res=>{
        localStorage.setItem("x",res.x.toString())
        localStorage.setItem("y",res.y.toString())
      });break;
      default:console.log("click btn_sys error");break;
    }
  }

  //预览输入按键
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

  //添加按键
  const addKey = ()=>{
    if(preKey.value==""){
      return;
    }
    if(vk[preKey.value]==undefined){
      ElMessage.warning({
        message: preKey.value+"键暂不支持！",
        offset: 50
      })
      return;
    }
    for (let i = 0; i < keyList.value.length; i++) {
      if(keyList.value[i].name == preKey.value){
        ElMessage.warning({
          message: preKey.value+"键冲突！",
          offset: 50
        })
        return ;
      }
    }
    keyList.value.push({name:preKey.value,value:vk[preKey.value],used:true,key_ms:"0"})
    onSaveKey()
  }

  //过滤快捷键列表
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
        if(parse.value == vq[key]){
          disabled = true
        }
        keyOptions.value.push({label:key,value:vq[key],disabled:disabled})
      }
    }
  }

  //保存快捷键
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

  //保存按键区键符及延迟值
  const onSaveKey=async () => {
    await initData()
    localStorage.setItem("keys",JSON.stringify(keyList.value))
  }

  //修改延迟值事件
  const changeMs = (val: { name: any;key_ms:number })=>{
    ElMessageBox.prompt(`当前值：${val.key_ms} (连发模式生效)`,`键[${val.name}]的间隔时间`, {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern:
          /^(5[0-9]|[6-9][0-9]|[1-9][0-9]{2,})$/,
      inputErrorMessage: '间隔时间必须为数字且>50',
    })
        .then(async ({value}) => {
          for (let i = 0; i < keyList.value.length; i++) {
            if (keyList.value[i].name == val.name) {
              keyList.value[i].key_ms = value
            }
          }
          await initData()
          localStorage.setItem("keys", JSON.stringify(keyList.value))
        })
        .catch(() => {})
  }

  //删除键符事件
  const onDelKey=async (value: { name: any; }) => {
    for (let i = 0; i < keyList.value.length; i++) {
      if (keyList.value[i].name == value.name) {
        keyList.value.splice(i, 1)
      }
    }
    await initData()
    localStorage.setItem("keys", JSON.stringify(keyList.value))
  }

  //初始化参数封装函数
  const initParam = (key: string, value: string) =>{
    let keyStr =  localStorage.getItem(key)
    if(keyStr == undefined){
      keyStr = value
      localStorage.setItem(key,keyStr)
    }
    return keyStr
  }

  //初始化键符
  const initData = async () => {
    let temp = []
    for (let i = 0; i < keyList.value.length; i++) {
      if (Boolean(keyList.value[i].used)) {
        temp.push({key:keyList.value[i].value,key_ms:keyList.value[i].key_ms})
      }
    }
    await SyncFrontKey(JSON.stringify(temp))
  }

  //清除快捷键及键符配置
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

  //开始结束音乐播放
  const playMusic=(isStart:boolean)=>{
    if(!voice.value){
      return
    }
    const start = document.getElementById('startMp3') as HTMLAudioElement; // 强制转换为 HTMLAudioElement
    const stop = document.getElementById('stopMp3') as HTMLAudioElement; // 强制转换为 HTMLAudioElement
    if(isStart){
      start.currentTime = 0;
      stop.pause()
      start.play()
    } else {
      stop.currentTime = 0;
      start.pause()
      stop.play()
    }
  }
  //初始化驱动文件
  const DllInit = async () => {
    let res = await DllImport()
    res == "success" ? LogDebug("驱动注入成功") : ElMessage.warning(res);
  }

  //初始化快捷键及键符
  const initKeys=()=>{
    start.value = initParam("start",vq['F9'])
    stop.value = initParam("stop",vq['F10'])
    parse.value = initParam("parse",vq['F11'])
    model.value = Number(initParam("model","0"))
    ms.value = Number(initParam("ms","50"))
    voice.value = initParam("voice","true")=="true"
    let keys =  localStorage.getItem("keys")
    if(keys == undefined){
      keyList.value = []
    } else {
      keyList.value = JSON.parse(keys)
    }
    filterKey(true)
  }



  //同步快捷键及键符到后台
  const syncKeys=async () => {
    await SyncFrontStart(Number(start.value))
    await SyncFrontStop(Number(stop.value))
    await SyncFrontParse(Number(parse.value))
    await SyncFrontMs(ms.value)
    await SyncFrontModel(model.value)
  }

  onMounted(async ()=>{
    const loading = ElLoading.service({
      lock: true,
      text: '正在加载，请稍后',
      background: 'rgba(255, 255, 255, 0.9)',
    })
    initKeys();
    await initData()
    await DllInit();
    await syncKeys()
    EventsOn("start-thread",()=>{
      isRunning.value = true
      playMusic(true)
    })
    EventsOn("stop-thread",()=>{
      isRunning.value = false
      playMusic(false)
    })
    loading.close()
  })
</script>
<style lang="less" scoped>
  @import "../../assets/style/theme";
  *{
    user-select: none;
  }
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