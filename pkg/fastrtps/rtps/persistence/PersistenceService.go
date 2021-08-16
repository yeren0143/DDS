package persistence

import (
	"dds/common"
	"dds/fastrtps/rtps/attributes"
	"dds/fastrtps/rtps/history"
)

// Abstract interface representing a persistence service implementaion
type IPersistenceService interface {
	/**
	 * Get all data stored for a writer.
	 * @param persistence_guid   GUID of the writer used to store samples.
	 * @param writer_guid        GUID of the writer to load.
	 * @param changes            History of the writer to load.
	 * @param change_pool        Pool where new changes should be obtained from.
	 * @param payload_pool       Pool where payloads should be obtained from.
	 * @param next_sequence      Sequence that should be applied to the next created sample.
	 * @return True if operation was successful.
	 */
	LoadWriterFromStoreage(persistenceGuID string, writerGUID *common.GUIDT, changes []*common.CacheChangeT,
		changePool history.IChangePool, payloadPool history.IPayloadPool,
		nextSequence *common.SequenceNumberT) bool

	/**
	 * Add a change to storage.
	 * @param persistence_guid   GUID of the writer used to store samples.
	 * @param change             The cache change to add.
	 * @return True if operation was successful.
	 */
	AddWriterChangeToStorage(persistenceGUID string, change *common.CacheChangeT) bool

	/**
	 * Remove a change from storage.
	 * @param persistence_guid   GUID of the writer used to store samples.
	 * @param change             The cache change to remove.
	 * @return True if operation was successful.
	 */
	RemoveWriterChangeFromStorage(persistenceGUID string, change *common.CacheChangeT) bool

	/**
	 * Get all data stored for a reader.
	 * @param reader_guid   GUID of the reader to load.
	 * @param seq_map       History record (map of low marks) to be loaded.
	 * @return True if operation was successful.
	 */
	LoadReaderFromStorage(readerGUID string, guid *common.GUIDT, seq *common.SequenceNumberT) bool

	/**
	 * Update the sequence number associated to a writer on a reader.
	 * @param reader_guid GUID of the reader to update.
	 * @param writer_guid GUID of the associated writer to update.
	 * @param seq_number New sequence number value to set for the associated writer.
	 * @return True if operation was successful.
	 */
	UpdateWriterSeqOnStorage(readerGUID string, writerGUID *common.GUIDT, seq *common.SequenceNumberT) bool
}

/**
 * Abstract factory to create a persistence service from participant or endpoint properties
 * @ingroup RTPS_PERSISTENCE_MODULE
 */
type IPersistenceFactory interface {
	/**
	* Create a persistence service implementation
	* @param property_policy PropertyPolicy where the persistence configuration will be searched
	* @return A pointer to a persistence service implementation. nullptr when policy does not contain the necessary properties or if persistence service could not be created
	 */
	CreatePersistenceService(*attributes.PropertyPolicyT) IPersistenceService
}
