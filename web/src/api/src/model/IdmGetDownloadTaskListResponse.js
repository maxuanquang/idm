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

import ApiClient from '../ApiClient';
import IdmDownloadTask from './IdmDownloadTask';

/**
 * The IdmGetDownloadTaskListResponse model module.
 * @module model/IdmGetDownloadTaskListResponse
 * @version version not set
 */
class IdmGetDownloadTaskListResponse {
    /**
     * Constructs a new <code>IdmGetDownloadTaskListResponse</code>.
     * @alias module:model/IdmGetDownloadTaskListResponse
     */
    constructor() { 
        
        IdmGetDownloadTaskListResponse.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>IdmGetDownloadTaskListResponse</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/IdmGetDownloadTaskListResponse} obj Optional instance to populate.
     * @return {module:model/IdmGetDownloadTaskListResponse} The populated <code>IdmGetDownloadTaskListResponse</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new IdmGetDownloadTaskListResponse();

            if (data.hasOwnProperty('downloadTaskList')) {
                obj['downloadTaskList'] = ApiClient.convertToType(data['downloadTaskList'], [IdmDownloadTask]);
            }
            if (data.hasOwnProperty('totalDownloadTaskCount')) {
                obj['totalDownloadTaskCount'] = ApiClient.convertToType(data['totalDownloadTaskCount'], 'String');
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>IdmGetDownloadTaskListResponse</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>IdmGetDownloadTaskListResponse</code>.
     */
    static validateJSON(data) {
        if (data['downloadTaskList']) { // data not null
            // ensure the json data is an array
            if (!Array.isArray(data['downloadTaskList'])) {
                throw new Error("Expected the field `downloadTaskList` to be an array in the JSON data but got " + data['downloadTaskList']);
            }
            // validate the optional field `downloadTaskList` (array)
            for (const item of data['downloadTaskList']) {
                IdmDownloadTask.validateJSON(item);
            };
        }
        // ensure the json data is a string
        if (data['totalDownloadTaskCount'] && !(typeof data['totalDownloadTaskCount'] === 'string' || data['totalDownloadTaskCount'] instanceof String)) {
            throw new Error("Expected the field `totalDownloadTaskCount` to be a primitive type in the JSON string but got " + data['totalDownloadTaskCount']);
        }

        return true;
    }


}



/**
 * @member {Array.<module:model/IdmDownloadTask>} downloadTaskList
 */
IdmGetDownloadTaskListResponse.prototype['downloadTaskList'] = undefined;

/**
 * @member {String} totalDownloadTaskCount
 */
IdmGetDownloadTaskListResponse.prototype['totalDownloadTaskCount'] = undefined;






export default IdmGetDownloadTaskListResponse;

