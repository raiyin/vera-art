import type { NewsDesc, } from '~/types';
import axios from 'axios';
import { useRuntimeConfig, } from '#imports';

const fetchCurrentNews = async (
    path: string,
): Promise<NewsDesc> => {
    try {
        const config = useRuntimeConfig();
        const newsid = path.substring(path.lastIndexOf('/',) + 1,);
        const response = await axios.get(config.public.serverUrl + 'news/' + newsid,);
        const oneCurrentNews = response.data;

        return oneCurrentNews;
    } catch (e) {
        console.log(e,);
        return {} as NewsDesc;
    }
};

const fetchOtherNews = async (
    path: string,
): Promise<NewsDesc[]> => {
    try {
        const config = useRuntimeConfig();
        const newsid = path.substring(path.lastIndexOf('/',) + 1,);
        const response = await axios.get(config.public.serverUrl + 'news', {
            params: { id_ne: newsid, limit: 5, },
        },);
        const otherNews = response.data.news ?? [];
        return otherNews;
    } catch (e) {
        console.log(e,);
        return [];
    }
};

export { fetchCurrentNews, fetchOtherNews, };
