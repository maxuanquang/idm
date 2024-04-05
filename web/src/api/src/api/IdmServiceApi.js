/**
 * idm.proto
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: version not set
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */


import ApiClient from "../ApiClient";
import IdmCreateAccountRequest from '../model/IdmCreateAccountRequest';
import IdmCreateAccountResponse from '../model/IdmCreateAccountResponse';
import IdmCreateDownloadTaskRequest from '../model/IdmCreateDownloadTaskRequest';
import IdmCreateDownloadTaskResponse from '../model/IdmCreateDownloadTaskResponse';
import IdmCreateSessionRequest from '../model/IdmCreateSessionRequest';
import IdmCreateSessionResponse from '../model/IdmCreateSessionResponse';
import IdmGetDownloadTaskListResponse from '../model/IdmGetDownloadTaskListResponse';
import IdmServiceUpdateDownloadTaskBody from '../model/IdmServiceUpdateDownloadTaskBody';
import IdmUpdateDownloadTaskResponse from '../model/IdmUpdateDownloadTaskResponse';
import RpcStatus from '../model/RpcStatus';
import StreamResultOfIdmGetDownloadTaskFileResponse from '../model/StreamResultOfIdmGetDownloadTaskFileResponse';

/**
* IdmService service.
* @module api/IdmServiceApi
* @version version not set
*/
export default class IdmServiceApi {

    /**
    * Constructs a new IdmServiceApi. 
    * @alias module:api/IdmServiceApi
    * @class
    * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
    * default to {@link module:ApiClient#instance} if unspecified.
    */
    constructor(apiClient) {
        this.apiClient = apiClient || ApiClient.instance;
    }


    /**
     * Callback function to receive the result of the idmServiceCreateAccount operation.
     * @callback module:api/IdmServiceApi~idmServiceCreateAccountCallback
     * @param {String} error Error message, if any.
     * @param {module:model/IdmCreateAccountResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {module:model/IdmCreateAccountRequest} body 
     * @param {module:api/IdmServiceApi~idmServiceCreateAccountCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/IdmCreateAccountResponse}
     */
    idmServiceCreateAccount(body, callback) {
      let postBody = body;
      // verify the required parameter 'body' is set
      if (body === undefined || body === null) {
        throw new Error("Missing the required parameter 'body' when calling idmServiceCreateAccount");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = IdmCreateAccountResponse;
      return this.apiClient.callApi(
        '/api/v1/accounts', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceCreateDownloadTask operation.
     * @callback module:api/IdmServiceApi~idmServiceCreateDownloadTaskCallback
     * @param {String} error Error message, if any.
     * @param {module:model/IdmCreateDownloadTaskResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {module:model/IdmCreateDownloadTaskRequest} body 
     * @param {module:api/IdmServiceApi~idmServiceCreateDownloadTaskCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/IdmCreateDownloadTaskResponse}
     */
    idmServiceCreateDownloadTask(body, callback) {
      let postBody = body;
      // verify the required parameter 'body' is set
      if (body === undefined || body === null) {
        throw new Error("Missing the required parameter 'body' when calling idmServiceCreateDownloadTask");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = IdmCreateDownloadTaskResponse;
      return this.apiClient.callApi(
        '/api/v1/tasks', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceCreateSession operation.
     * @callback module:api/IdmServiceApi~idmServiceCreateSessionCallback
     * @param {String} error Error message, if any.
     * @param {module:model/IdmCreateSessionResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {module:model/IdmCreateSessionRequest} body 
     * @param {module:api/IdmServiceApi~idmServiceCreateSessionCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/IdmCreateSessionResponse}
     */
    idmServiceCreateSession(body, callback) {
      let postBody = body;
      // verify the required parameter 'body' is set
      if (body === undefined || body === null) {
        throw new Error("Missing the required parameter 'body' when calling idmServiceCreateSession");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = IdmCreateSessionResponse;
      return this.apiClient.callApi(
        '/api/v1/sessions', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceDeleteDownloadTask operation.
     * @callback module:api/IdmServiceApi~idmServiceDeleteDownloadTaskCallback
     * @param {String} error Error message, if any.
     * @param {Object} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {String} downloadTaskId 
     * @param {module:api/IdmServiceApi~idmServiceDeleteDownloadTaskCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Object}
     */
    idmServiceDeleteDownloadTask(downloadTaskId, callback) {
      let postBody = null;
      // verify the required parameter 'downloadTaskId' is set
      if (downloadTaskId === undefined || downloadTaskId === null) {
        throw new Error("Missing the required parameter 'downloadTaskId' when calling idmServiceDeleteDownloadTask");
      }

      let pathParams = {
        'downloadTaskId': downloadTaskId
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = Object;
      return this.apiClient.callApi(
        '/api/v1/tasks/{downloadTaskId}', 'DELETE',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceDeleteSession operation.
     * @callback module:api/IdmServiceApi~idmServiceDeleteSessionCallback
     * @param {String} error Error message, if any.
     * @param {Object} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {module:api/IdmServiceApi~idmServiceDeleteSessionCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Object}
     */
    idmServiceDeleteSession(callback) {
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = Object;
      return this.apiClient.callApi(
        '/api/v1/sessions', 'DELETE',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceGetDownloadTaskFile operation.
     * @callback module:api/IdmServiceApi~idmServiceGetDownloadTaskFileCallback
     * @param {String} error Error message, if any.
     * @param {module:model/StreamResultOfIdmGetDownloadTaskFileResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {String} downloadTaskId 
     * @param {module:api/IdmServiceApi~idmServiceGetDownloadTaskFileCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/StreamResultOfIdmGetDownloadTaskFileResponse}
     */
    idmServiceGetDownloadTaskFile(downloadTaskId, callback) {
      let postBody = null;
      // verify the required parameter 'downloadTaskId' is set
      if (downloadTaskId === undefined || downloadTaskId === null) {
        throw new Error("Missing the required parameter 'downloadTaskId' when calling idmServiceGetDownloadTaskFile");
      }

      let pathParams = {
        'downloadTaskId': downloadTaskId
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = StreamResultOfIdmGetDownloadTaskFileResponse;
      return this.apiClient.callApi(
        '/api/v1/tasks/{downloadTaskId}/files', 'GET',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceGetDownloadTaskList operation.
     * @callback module:api/IdmServiceApi~idmServiceGetDownloadTaskListCallback
     * @param {String} error Error message, if any.
     * @param {module:model/IdmGetDownloadTaskListResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {Object} opts Optional parameters
     * @param {String} [offset] 
     * @param {String} [limit] 
     * @param {module:api/IdmServiceApi~idmServiceGetDownloadTaskListCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/IdmGetDownloadTaskListResponse}
     */
    idmServiceGetDownloadTaskList(opts, callback) {
      opts = opts || {};
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
        'offset': opts['offset'],
        'limit': opts['limit']
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = IdmGetDownloadTaskListResponse;
      return this.apiClient.callApi(
        '/api/v1/tasks', 'GET',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the idmServiceUpdateDownloadTask operation.
     * @callback module:api/IdmServiceApi~idmServiceUpdateDownloadTaskCallback
     * @param {String} error Error message, if any.
     * @param {module:model/IdmUpdateDownloadTaskResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * @param {String} downloadTaskId 
     * @param {module:model/IdmServiceUpdateDownloadTaskBody} body 
     * @param {module:api/IdmServiceApi~idmServiceUpdateDownloadTaskCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/IdmUpdateDownloadTaskResponse}
     */
    idmServiceUpdateDownloadTask(downloadTaskId, body, callback) {
      let postBody = body;
      // verify the required parameter 'downloadTaskId' is set
      if (downloadTaskId === undefined || downloadTaskId === null) {
        throw new Error("Missing the required parameter 'downloadTaskId' when calling idmServiceUpdateDownloadTask");
      }
      // verify the required parameter 'body' is set
      if (body === undefined || body === null) {
        throw new Error("Missing the required parameter 'body' when calling idmServiceUpdateDownloadTask");
      }

      let pathParams = {
        'downloadTaskId': downloadTaskId
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = IdmUpdateDownloadTaskResponse;
      return this.apiClient.callApi(
        '/api/v1/tasks/{downloadTaskId}', 'PUT',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }


}
