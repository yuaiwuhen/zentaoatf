import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { QueryParams, QueryResult } from '@/types/data.d';
import {
    list, get, remove,getLastest,
} from './service';

export interface StateType {
    queryResult: QueryResult;
    detailResult: any;
    lastResult: any;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setQueryResult: Mutation<StateType>;
        setDetailResult: Mutation<StateType>;
        setLastResult: Mutation<StateType>;
    };
    actions: {
        list: Action<StateType, StateType>;
        get: Action<StateType, StateType>;
        delete: Action<StateType, StateType>;
        latest: Action<StateType, StateType>;
    };
}
const initState: StateType = {
    queryResult: {
        result: [],
        pagination: {
            total: 0,
            page: 1,
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
        },
    },
    detailResult: {},
    lastResult: {},
};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'Result',
    state: {
        ...initState
    },
    mutations: {
        setQueryResult(state, payload) {
            state.queryResult = payload;
        },
        setDetailResult(state, payload) {
            state.detailResult = payload;
        },
        setLastResult(state, payload) {
            state.lastResult = payload;
        },
    },
    actions: {
        async list({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await list(params);
                if (response.code != 0) {
                    return;
                }
                const data = response.data;
                commit('setQueryResult', data);

                return true;
            } catch (error) {
                return false;
            }
        },
        async get({ commit }, params: any ) {
            const response: ResponseData = await get(params);
            const data = response.data;

            commit('setDetailResult',data);

            return true;
        },
        async delete({ commit , dispatch, state}, data: any ) {
            try {
                await remove(data);
                dispatch('list', {})

                return true;
            } catch (error) {
                return false;
            }
        },
        async latest({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await getLastest(params);
                if (response.code != 0) {
                    return;
                }
                const data = response.data;
                commit('setLastResult', data);

                return true;
            } catch (error) {
                return false;
            }
        },
    }
};

export default StoreModel;
