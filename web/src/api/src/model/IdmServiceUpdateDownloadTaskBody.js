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
import IdmDownloadStatus from './IdmDownloadStatus';

/**
 * The IdmServiceUpdateDownloadTaskBody model module.
 * @module model/IdmServiceUpdateDownloadTaskBody
 * @version version not set
 */
class IdmServiceUpdateDownloadTaskBody {
    /**
     * Constructs a new <code>IdmServiceUpdateDownloadTaskBody</code>.
     * @alias module:model/IdmServiceUpdateDownloadTaskBody
     */
    constructor() { 
        
        IdmServiceUpdateDownloadTaskBody.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>IdmServiceUpdateDownloadTaskBody</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/IdmServiceUpdateDownloadTaskBody} obj Optional instance to populate.
     * @return {module:model/IdmServiceUpdateDownloadTaskBody} The populated <code>IdmServiceUpdateDownloadTaskBody</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new IdmServiceUpdateDownloadTaskBody();

            if (data.hasOwnProperty('downloadStatus')) {
                obj['downloadStatus'] = IdmDownloadStatus.constructFromObject(data['downloadStatus']);
            }
            if (data.hasOwnProperty('metadata')) {
                obj['metadata'] = ApiClient.convertToType(data['metadata'], 'String');
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>IdmServiceUpdateDownloadTaskBody</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>IdmServiceUpdateDownloadTaskBody</code>.
     */
    static validateJSON(data) {
        // ensure the json data is a string
        if (data['metadata'] && !(typeof data['metadata'] === 'string' || data['metadata'] instanceof String)) {
            throw new Error("Expected the field `metadata` to be a primitive type in the JSON string but got " + data['metadata']);
        }

        return true;
    }


}



/**
 * @member {module:model/IdmDownloadStatus} downloadStatus
 */
IdmServiceUpdateDownloadTaskBody.prototype['downloadStatus'] = undefined;

/**
 * @member {String} metadata
 */
IdmServiceUpdateDownloadTaskBody.prototype['metadata'] = undefined;






export default IdmServiceUpdateDownloadTaskBody;

