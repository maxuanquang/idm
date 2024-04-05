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
* Enum class IdmDownloadType.
* @enum {}
* @readonly
*/
export default class IdmDownloadType {
    
        /**
         * value: "UndefinedType"
         * @const
         */
        "UndefinedType" = "UndefinedType";

    
        /**
         * value: "HTTP"
         * @const
         */
        "HTTP" = "HTTP";

    

    /**
    * Returns a <code>IdmDownloadType</code> enum value from a Javascript object name.
    * @param {Object} data The plain JavaScript object containing the name of the enum value.
    * @return {module:model/IdmDownloadType} The enum <code>IdmDownloadType</code> value.
    */
    static constructFromObject(object) {
        return object;
    }
}
