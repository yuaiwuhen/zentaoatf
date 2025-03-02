<template>
  <div id="log-list" class="log-list scrollbar-y">
    <pre id="content" class="content">
      <template v-for="(item, index) in wsMsg.out" :key="index">
        {{ void (info = item.info) }}
        {{ void (csKey = info?.key) }}

        <code class="item small"
             :class="[
                 csKey && caseDetail[csKey] ? 'show-detail' : '',

                 csKey ? 'case-item' : '',
                 info?.status === 'start' ? 'case-start' : '',
                 info?.status === 'start' ? 'result-'+caseResult[csKey] : '',

                 info?.status === 'start-task' ? 'strong' : ''
             ]">

          <div class="group">
            <template v-if="info?.status === 'start'">
              <span @click="showDetail(item.info?.key)" class="link state center">
                <Icon v-if="!caseDetail[csKey]" icon="chevron-right" />
                <Icon v-if="caseDetail[csKey]" icon="chevron-down" />
              </span>
            </template>
          </div>

          <div class="sign">
            <Icon v-if="item.msg" icon="circle" />
            <span v-else>&nbsp;</span>
          </div>

          <div class="time">
            <span>{{ item.time }}</span>
          </div>
          <div class="msg-span">
            <span v-html="item.msg"></span>
            <span v-if="info?.status === 'start' && caseResult[csKey]">
              [ {{ t(caseResult[csKey]) }} ]
            </span>
          </div>
        </code>

      </template>
    </pre>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {WsMsg} from "@/types/data";
import {genExecInfo, genWorkspaceToScriptsMap} from "@/views/script/service";
import {updateStatistic} from "@/views/result/service";
import {scroll} from "@/utils/dom";
import {computed, onBeforeUnmount, onMounted, reactive, ref, watch} from "vue";
import {useStore} from "vuex";
import {WebSocketData} from "@/store/websoket";
import {ExecStatus} from "@/store/exec";
import {ProxyData} from "@/store/proxy";
import {WorkspaceData} from "@/store/workspace";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {ZentaoData} from "@/store/zentao";
import {WebSocket} from "@/services/websocket";

import Icon from '@/components/Icon.vue';
import {momentTime} from "@/utils/datetime";
import {isInArray} from "@/utils/array";
import {StateType as GlobalStateType} from "@/store/global";
import { get as getWorkspace, uploadToProxy } from "@/views/workspace/service";
import {mvLog} from "@/views/result/service";
import {checkProxy} from "@/views/proxy/service";
import useCurrentProduct from "@/hooks/use-current-product";

const { t } = useI18n();

const store = useStore<{global: GlobalStateType, Zentao: ZentaoData, WebSocket: WebSocketData, Exec: ExecStatus, workspace: WorkspaceData, proxy: ProxyData}>();
const logContentExpand = computed<boolean>(() => store.state.global.logContentExpand);

const currProduct = computed<any>(() => store.state.Zentao.currProduct);
const isRunning = computed<any>(() => store.state.Exec.isRunning);
const currentWorkspace = ref({} as any);
const defaultProxy = computed<any>(() => store.state.proxy.currProxy);
const currProxy = ref({} as any);

const cachedExecData = ref({})
const caseCount = ref(1)
const caseResult = ref({})
const caseDetail = ref({})
const realPathMap = ref({})
const lastWsMsg = ref('')

// websocket
let wsMsg = reactive({in: '', out: [] as any[]});

const showDetail = (key) => {
  console.log('showDetail', key)
  caseDetail.value[key] = !caseDetail.value[key]
}

const hideAllDetailOrNot = (val) => {
  console.log('hideAllDetailOrNot')
  Object.keys(caseDetail.value).forEach((key => {
    caseDetail.value[key] = val;
  }))
  console.log(caseDetail.value)
}

watch(logContentExpand, () => {
  console.log('watch logContentExpand', logContentExpand.value)
  hideAllDetailOrNot(logContentExpand.value)
}, {deep: true})

const onWebsocketMsgEvent = async (data: any) => {
  console.log('WebsocketMsgEvent in ExecLog', data.msg)

  let item = JSON.parse(data.msg) as WsMsg
  if(item.category == 'watch'){
    return;
  }
  if (item.isRunning != undefined) {
    console.log('change isRunning to ', item.isRunning)
    store.dispatch('Exec/setRunning', item.isRunning)
  }

  if (item.info?.status === 'start') {
    const key = item.info.key + '-' + caseCount.value
    caseDetail.value[key] = logContentExpand.value
  }

  item = genExecInfo(item, caseCount.value)

  if(lastWsMsg.value === JSON.stringify(item)){
    return;
  }
  lastWsMsg.value = JSON.stringify(item);

  if (item.info && item.info.key && isInArray(item.info.status, ['pass', 'fail', 'skip'])) { // set case result
    store.dispatch('Result/list', {
        keywords: '',
        enabled: 1,
        pageSize: 10,
        page: 1
        });
    caseResult.value[item.info.key] = item.info.status
  }

  if(currProxy.value.path != undefined && currProxy.value.path != '' && currProxy.value.path != 'local'){
    if(item.info != undefined){
        item.info.logDir = await downloadLog(item);
    }
    if( item.msg == ""){
      return;
    }
  }

  Object.keys(realPathMap.value).forEach(key => {
    item.msg = item.msg.replace(key, realPathMap.value[key])
  })
  
  if(item.info?.logDir != undefined && currentWorkspace.value.type == 'ztf'){
    updateStatisticInfo(item.info.logDir)
  }
  wsMsg.out.push(item)
  scroll('log-list')
}

const downloadLog = async (item) => {
  const nItem = {...item}
  const msg = nItem.msg;
  let pth = nItem?.logDir
  if(msg.indexOf('Report') !== 0 && msg.indexOf('报告') !== 0){
    return;
  }
  let logPath = '';
  if(msg.indexOf('Report') === 0){
    logPath = msg.replace('Report', '')
    logPath = logPath.replace('.', '')
  }else if(msg.indexOf('报告') === 0){
    logPath = msg.replace('报告', '')
    logPath = logPath.replace('。', '')
  }

  await mvLog({file: logPath, workspaceId: currentWorkspace.value.id, proxyPath: currProxy.value.path, pathMap: realPathMap.value}).then(resp => {
    if(resp.code === 0){
      pth = resp.data.replace('log.txt', '');
      item.msg = nItem.msg.replace(logPath, resp.data)
      let emptMsg = {...nItem}
      emptMsg.msg = '';
      emptMsg.isRunning = false;
      emptMsg.time = '';
      setTimeout(() => {
        wsMsg.out.push(emptMsg)
      }, 100);
      store.dispatch('Result/list', {
        keywords: '',
        enabled: 1,
        pageSize: 10,
        page: 1
        });
    }
  })
  return pth
}
const updateStatisticInfo = (logDir) => {
  console.log('updateStatisticInfo')
  updateStatistic({path: logDir}).then((res) => {
    if(res.code == 0){
        res.data.forEach((item) => {
            store.dispatch('Result/getStatistic', {path: item})
        })
    }
  })
}

const clearWebsocketMsgEvent = () => {
  wsMsg.out = [];
}

onMounted(() => {
  console.log('onMounted')
  bus.on(settings.eventExec, exec);
  bus.on(settings.eventWebSocketMsg, onWebsocketMsgEvent);
  bus.on(settings.eventClearWebSocketMsg, clearWebsocketMsgEvent);
})
onBeforeUnmount( () => {
  console.log('onBeforeUnmount')
  bus.off(settings.eventExec, exec);
  bus.off(settings.eventWebSocketMsg, onWebsocketMsgEvent);
  bus.off(settings.eventClearWebSocketMsg, clearWebsocketMsgEvent);
})

const exec = async (data: any) => {
  console.log('exec', data)

  let execType = data.execType
  if (execType === 'previous') {
    data = cachedExecData.value
    execType = data.execType
  } else {
    cachedExecData.value = data
  }

  if (execType === 'ztf' && (!data.scripts || data.scripts.length === 0)) {
    const msgCancel = {
      msg: `<span class="strong">`+t('case_num_empty')+`</span>`,
      time: momentTime(new Date())}
    wsMsg.out.push(msgCancel)
    return
  }

  caseCount.value++
  let msg = {} as any
  let workspaceId = 0;
  if (execType === 'ztf') {
    const scripts = data.scripts
    const sets = genWorkspaceToScriptsMap(scripts)
    if(!workspaceId){
        if(scripts.length > 0){
            workspaceId = scripts[0].workspaceId
        }
    }
    if(!workspaceId){
        if(scripts.length > 0){
            workspaceId = scripts[0].workspaceId
        }
    }
    console.log('===', sets)
    msg = {act: 'execCase', testSets: sets}

  } else if (execType === 'unit') {
    const set = {
      workspaceId: data.id,
      workspaceType: data.type,
      cmd: data.cmd,
      submitResult: data.submitResult,
      name: data.name,
    }
    workspaceId = workspaceId == 0 ? data.id : workspaceId;

    msg = {act: 'execUnit', testSets: [set], productId: currProduct.value.id}

  } else if (execType === 'stop') {
    msg = {act: 'execStop'}
  }

  console.log('exec testing', msg)

  currProxy.value = await checkProxyStatus(workspaceId, msg)
  if(currProxy.value.path == ''){
      return;
  }
  await uploadScript(workspaceId, msg)
  WebSocket.sentMsg(settings.webSocketRoom, JSON.stringify(msg), currProxy.value.path)
}

const checkProxyStatus = async (workspaceId, msg) => {
  if(msg.testSets == undefined){
    return {path: 'local'};
  }
  currentWorkspace.value = {};
  let workspaceInfo = {} as any;
  if(workspaceId > 0){
    workspaceInfo = await getWorkspace(workspaceId)
  }
  const selectedProxy = defaultProxy.value
  if(selectedProxy.id > 0){
    const msgSelectProxy = {
        msg: `<span class="strong">`+t('case_test_proxy')+`</span>`,
        time: momentTime(new Date())}
    wsMsg.out.push(msgSelectProxy)
    const resp = await checkProxy(selectedProxy.id)
    if(resp.data.status != 'ok'){
      const msgSelectProxy = {
        msg: `<span class="strong result-fail">`+t('biz_5002')+`</span>`,
        time: momentTime(new Date())
      }
      wsMsg.out.push(msgSelectProxy)
      return {path: ''};
    }
  }
  if(workspaceInfo.data.proxy_id != selectedProxy.id){
    workspaceInfo.data.proxy_id = selectedProxy.id
    console.log("change workspace proxyid")
    store.dispatch('Workspace/save', workspaceInfo.data)
  }
  currentWorkspace.value = workspaceInfo.data;
  return selectedProxy;
}

const uploadScript = async (workspaceId, msg) => {
  if(msg.testSets == undefined || currProxy.value.id == 0 || currProxy.value.id == undefined){
    return '';
  }
  const msgUpload = {
  msg: `<span class="strong">`+t('case_upload_to_proxy')+`</span>`,
  time: momentTime(new Date())}
  wsMsg.out.push(msgUpload)
  msg.testSets[0].proxyId = currProxy.value.id;
  const resp = await uploadToProxy(msg.testSets);
  const testSetsMap = resp.data;
  realPathMap.value = testSetsMap;
  let casesMap = {};
  var keys = Object.keys(testSetsMap);
  keys.forEach((val) => {
    casesMap[testSetsMap[val]] = val;
  });
  msg.testSets.forEach((set, setIndex) => {
    if (set.cases != undefined) {
      set.cases.forEach((casePath, caseIndex) => {
      msg.testSets[setIndex].cases[caseIndex] = casesMap[casePath];
      });
    }else{
      msg.testSets[setIndex].workspacePath = keys[0]
      msg.testSets[setIndex].workspaceType = currentWorkspace.value.type
    }
  });
}

</script>

<style lang="less">
.log-list {
  height: 100%;
  .result-pass {
    color: var(--color-green)
  }
  .result-fail {
    color: var(--color-red)
  }
  .result-skip {
    color: var(--color-gray)
  }

  .content {
    white-space: normal;
    .item {
      line-height: 18px;
      &.case-item {
        &:not(.case-start) {
          display: none !important;
        }
      }
      &.show-detail:not(.case-start) {
        display: flex !important;
      }
      &:hover {
        background-color: var(--color-darken-1);
      }

      .group {
        width: 16px;
        font-size: 13px;
        text-align: center;
        .link {
          //position: relative;
          //top: 2px;
        }
      }
      .sign {
        // margin: auto;

        width: 20px;
        font-size: 6px;
        text-align: center;
      }
      .time {
        width: 80px;
      }
    }
  }
}
</style>
