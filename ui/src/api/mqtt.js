import {useAxios} from '@/config/axios/axios';

const request = useAxios();

// 查询mqtt客户端列表
export const queryMqttClientList = params => {
    return request.get({ url: '/mqtt/client/list', params });
};

// 新增mqtt客户端
export const saveMqttClient = data => {
    return request.post({ url: '/mqtt/client/add', data });
}

// 删除mqtt客户端
export const delMqttClient = data => {
    return request.delete({ url: '/mqtt/client/delete', data });
}



// 转发器相关接口
export const saveForwarder = data => {
    return request.post({ url: '/mqtt/forwarder/add', data});
}

// 删除转发器
export const delForwarder = data => {
    return request.delete({ url: '/mqtt/forwarder/delete', data});
}

// 获取转发器列表
export const listForwarder = params => {
    return request.get({ url: '/mqtt/forwarder/list', params});
}

// 开启转发器
export const enableForwarder = id => {
    return request.post({ url: `/mqtt/forwarder/enable/${id}`  });
}

// 关闭转发器
export const disableForwarder = id => {
    return request.post({ url: `/mqtt/forwarder/disable/${id}` });
}

// ====================
// 获取聚英继电器
export const connectRelay = data => {
    return request.post({url: '/mqtt/relay/connect', data});
}

export const disconnectRelay = data => {
    return request.post({url: '/mqtt/relay/disconnect', data});
}


export const switchStatus = data => {
    return request.post({url: '/mqtt/relay/command', data});
}

export const listRelay = params => {
    return request.get({ url: '/mqtt/relay/list', params});
}