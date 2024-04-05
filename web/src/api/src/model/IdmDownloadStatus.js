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
/**
* Enum class IdmDownloadStatus.
* @enum {}
* @readonly
*/
export default class IdmDownloadStatus {
    
        /**
         * value: "UndefinedStatus"
         * @const
         */
        "UndefinedStatus" = "UndefinedStatus";

    
        /**
         * value: "Pending"
         * @const
         */
        "Pending" = "Pending";

    
        /**
         * value: "Downloading"
         * @const
         */
        "Downloading" = "Downloading";

    
        /**
         * value: "Failed"
         * @const
         */
        "Failed" = "Failed";

    
        /**
         * value: "Success"
         * @const
         */
        "Success" = "Success";

    

    /**
    * Returns a <code>IdmDownloadStatus</code> enum value from a Javascript object name.
    * @param {Object} data The plain JavaScript object containing the name of the enum value.
    * @return {module:model/IdmDownloadStatus} The enum <code>IdmDownloadStatus</code> value.
    */
    static constructFromObject(object) {
        return object;
    }
}
