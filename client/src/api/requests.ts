import type { NewsItemType } from '../types';
import axios from 'axios';

const fetchCurrentNews = async (
    path: string
): Promise<NewsItemType> => {
    try {
        const newsid = path.substring(path.lastIndexOf('/') + 1);
        const response = await axios.get(import.meta.env.VITE_SERVER_URL + 'news', {
            params: { id: newsid },
        });
        const oneCurrentNews = response.data[0];

        return oneCurrentNews;
    } catch (e) {
        console.log(e);
        return {} as NewsItemType;
    }
};

const fetchOtherNews = async (
    path: string
): Promise<NewsItemType[]> => {
    try {
        const newsid = path.substring(path.lastIndexOf('/') + 1);
        const response = await axios.get(import.meta.env.VITE_SERVER_URL + 'news', {
            params: { id_ne: newsid, limit: 5 },
        });
        const otherNews = response.data;
        return otherNews;
    } catch (e) {
        console.log(e);
        return [];
    }
};

export { fetchCurrentNews, fetchOtherNews };
