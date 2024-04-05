import { ApiClient, IdmServiceApi } from "../api/src";

const apiClient = new ApiClient();
apiClient.basePath = "http://localhost:8081";
apiClient.enableCookies = true;

const idmServiceApi = new IdmServiceApi(apiClient);

export default idmServiceApi;
