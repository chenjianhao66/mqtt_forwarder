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