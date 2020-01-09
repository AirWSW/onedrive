package graphapi

import (
	"time"
)

type NewMicrosoftGraphAPIInput struct {
	MicrosoftEndPoints     *MicrosoftEndPoints     `json:"microsoftEndPoints"`
	AzureADAppRegistration *AzureADAppRegistration `json:"azureAdAppRegistration"`
	AzureADAuthFlowContext *AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
}

type MicrosoftGraphAPI struct {
	MicrosoftEndPoints     MicrosoftEndPoints      `json:"microsoftEndPoints"`
	AzureADAppRegistration AzureADAppRegistration  `json:"azureAdAppRegistration"`
	AzureADAuthFlowContext AzureADAuthFlowContext  `json:"azureAdAuthFlowContext"`
	MicrosoftGraphAPIToken *MicrosoftGraphAPIToken `json:"microsoftGraphApiToken"`
}

type MicrosoftEndPoints struct {
	AzureADPortalEndPointURL     *string `json:"azureAdPortalEndPointUrl"`
	AzureADEndPointURL           string  `json:"azureAdEndPointUrl"`
	MicrosoftGraphAPIEndPointURL string  `json:"microsoftgraphApiEndPointUrl"`
}

type AzureADAppRegistration struct {
	DisplayName  *string  `json:"displayName"`
	ClientID     string   `json:"clientId"`
	TenantID     *string  `json:"tenantId"`
	ObjectID     *string  `json:"objectId"`
	RedirectURIs []string `json:"redirectUris"`
	LogoutURL    *string  `json:"logoutUrl"`
	ClientSecret string   `json:"clientSecret"`
}

type AzureADAuthFlowContext struct {
	GrantScope   string  `json:"grantScope"`
	Code         *string `json:"code"`
	RefreshToken *string `json:"refreshToken"`
}

type MicrosoftGraphAPIToken struct {
	TokenType    string  `json:"token_type"`
	ExpiresIn    int32   `json:"expires_in"`
	ExtExpiresIn *int32  `json:"ext_expires_in"`
	Scope        string  `json:"scope"`
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}

// MicrosoftGraphBaseItem  "@odata.type": "microsoft.graph.baseItem"
type MicrosoftGraphBaseItem struct {
	ID                   string                       `json:"id"` // identifier
	CreatedBy            *MicrosoftGraphIdentitySet   `json:"createdBy"`
	CreatedDateTime      *time.Time                   `json:"createdDateTime"`
	Description          string                       `json:"description"`
	ETag                 string                       `json:"eTag"`
	LastModifiedBy       *MicrosoftGraphIdentitySet   `json:"lastModifiedBy"`
	LastModifiedDateTime time.Time                    `json:"lastModifiedDateTime"`
	Name                 string                       `json:"name"`
	ParentReference      *MicrosoftGraphItemReference `json:"parentReference"`
	WebURL               string                       `json:"webUrl"`
}

// MicrosoftGraphDriveItem  "@odata.type": "microsoft.graph.driveItem"
type MicrosoftGraphDriveItem struct {
	Audio          *MicrosoftGraphAudio            `json:"audio"`
	Content        *EdmDotStream                   `json:"content"`
	CTag           string                          `json:"cTag"` // etag
	Deleted        *MicrosoftGraphDeleted          `json:"deleted"`
	Description    string                          `json:"description"`
	File           *MicrosoftGraphFile             `json:"file"`
	FileSystemInfo *MicrosoftGraphFileSystemInfo   `json:"fileSystemInfo"`
	Folder         *MicrosoftGraphFolder           `json:"folder"`
	Image          *MicrosoftGraphImage            `json:"image"`
	Location       *MicrosoftGraphGEOCoordinates   `json:"location"`
	Malware        *MicrosoftGraphMalware          `json:"malware"`
	Package        *MicrosoftGraphPackage          `json:"package"`
	Photo          *MicrosoftGraphPhoto            `json:"photo"`
	Publication    *MicrosoftGraphPublicationFacet `json:"publication"`
	RemoteItem     *MicrosoftGraphRemoteItem       `json:"remoteItem"`
	Root           *MicrosoftGraphRoot             `json:"root"`
	SearchResult   *MicrosoftGraphSearchResult     `json:"searchResult"`
	Shared         *MicrosoftGraphShared           `json:"shared"`
	SharepointIDs  *MicrosoftGraphSharepointIDs    `json:"sharepointIds"`
	Size           int64                           `json:"size"`
	SpecialFolder  *MicrosoftGraphSpecialFolder    `json:"specialFolder"`
	Video          *MicrosoftGraphVideo            `json:"video"`
	WebDavURL      string                          `json:"webDavUrl"`

	/* relationships */
	Activities  []MicrosoftGraphItemActivity     `json:"activities"`
	Children    []MicrosoftGraphDriveItem        `json:"children"`
	Permissions []MicrosoftGraphPermission       `json:"permissions"`
	Thumbnails  []MicrosoftGraphThumbnailSet     `json:"thumbnails"`
	Versions    []MicrosoftGraphDriveItemVersion `json:"versions"`

	/* inherited from baseItem */
	ID                   string                       `json:"id"` // identifier
	CreatedBy            *MicrosoftGraphIdentitySet   `json:"createdBy"`
	CreatedDateTime      *time.Time                   `json:"createdDateTime"`
	ETag                 string                       `json:"eTag"`
	LastModifiedBy       *MicrosoftGraphIdentitySet   `json:"lastModifiedBy"`
	LastModifiedDateTime *time.Time                   `json:"lastModifiedDateTime"`
	Name                 string                       `json:"name"`
	ParentReference      *MicrosoftGraphItemReference `json:"parentReference"`
	WebURL               string                       `json:"webUrl"`

	/* instance annotations */
	AtMicrosoftGraphConflictBehavior string `json:"@microsoft.graph.conflictBehavior"`
	AtMicrosoftGraphDownloadURL      string `json:"@microsoft.graph.downloadUrl"`
	AtMicrosoftGraphSourceURL        string `json:"@microsoft.graph.sourceUrl"`
}

// MicrosoftGraphDriveItemCollection "@odata.type": "microsoft.graph.driveItemCollection"
type MicrosoftGraphDriveItemCollection struct {
	Value           []MicrosoftGraphDriveItem `json:"value"`
	AtODataNextLink string                    `json:"@odata.nextLink"`
}

// MicrosoftGraphAudio "@odata.type": "microsoft.graph.audio"
type MicrosoftGraphAudio struct {
	Album             string `json:"album"`
	AlbumArtist       string `json:"albumArtist"`
	Artist            string `json:"artist"`
	Bitrate           int64  `json:"bitrate"`
	Composers         string `json:"composers"`
	Copyright         string `json:"copyright"`
	Disc              int16  `json:"disc"`
	DiscCount         int16  `json:"discCount"`
	Duration          int64  `json:"duration"`
	Genre             string `json:"genre"`
	HasDRM            bool   `json:"hasDrm"`
	IsVariableBitrate bool   `json:"isVariableBitrate"`
	Title             string `json:"title"`
	Track             int32  `json:"track"`
	TrackCount        int32  `json:"trackCount"`
	Year              int32  `json:"year"`
}

// EdmDotStream "@odata.type": "Edm.Stream"
type EdmDotStream struct {
	// The content stream.
}

// MicrosoftGraphIdentitySet "@odata.type": "microsoft.graph.identitySet"
type MicrosoftGraphIdentitySet struct {
	Application *MicrosoftGraphIdentity `json:"application"`
	Device      *MicrosoftGraphIdentity `json:"device"`
	Group       *MicrosoftGraphIdentity `json:"group"`
	User        *MicrosoftGraphIdentity `json:"user"`
}

// MicrosoftGraphIdentity "@odata.type": "microsoft.graph.identity"
type MicrosoftGraphIdentity struct {
	DisplayName string                      `json:"displayName"`
	ID          string                      `json:"id"`
	Thumbnails  *MicrosoftGraphThumbnailSet `json:"thumbnails"`
}

// MicrosoftGraphDeleted "@odata.type": "microsoft.graph.deleted"
type MicrosoftGraphDeleted struct {
	State string `json:"state"`
}

// MicrosoftGraphFile "@odata.type": "microsoft.graph.file"
type MicrosoftGraphFile struct {
	Hashes             string `json:"hashes"`
	MimeType           string `json:"mimeType"`
	ProcessingMetadata bool   `json:"processingMetadata"`
}

// MicrosoftGraphHashes "@odata.type": "microsoft.graph.hashes"
type MicrosoftGraphHashes struct {
	CRC32Hash    string `json:"crc32Hash"`    // hex
	SHA1Hash     string `json:"sha1Hash"`     // hex
	QuickXorHash string `json:"quickXorHash"` // base64
}

// MicrosoftGraphFileSystemInfo "@odata.type": "microsoft.graph.fileSystemInfo"
type MicrosoftGraphFileSystemInfo struct {
	CreatedDateTime      *time.Time `json:"createdDateTime"`
	LastAccessedDateTime *time.Time `json:"lastAccessedDateTime"`
	LastModifiedDateTime *time.Time `json:"lastModifiedDateTime"`
}

// MicrosoftGraphFolder "@odata.type": "microsoft.graph.folder"
type MicrosoftGraphFolder struct {
	ChildCount int32                     `json:"childCount"`
	View       *MicrosoftGraphFolderView `json:"view"`
}

// MicrosoftGraphFolderView "@odata.type": "microsoft.graph.folderView"
type MicrosoftGraphFolderView struct {
	SortBy    string `json:"sortBy"`    // default, name, type, size, takenOrCreatedDateTime, lastModifiedDateTime, sequence
	SortOrder string `json:"sortOrder"` // ascending, descending
	ViewType  string `json:"viewType"`  // default, icons, details, thumbnails
}

// MicrosoftGraphImage "@odata.type": "microsoft.graph.image"
type MicrosoftGraphImage struct {
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

// MicrosoftGraphGEOCoordinates "@odata.type": "microsoft.graph.geoCoordinates"
type MicrosoftGraphGEOCoordinates struct {
	Altitude  float64 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// MicrosoftGraphMalware "@odata.type": "microsoft.graph.malware"
type MicrosoftGraphMalware struct {
	/* The Malware resource indicates that malware has been detected in the item.
	 * In this version of the API, the presence (non-null) of the resource indicates
	 * that item contains malware, while a null (or missing) value indicates that
	 * the item is clean.
	 *
	 * Note: While this resource is empty today, in future API revisions the resource
	 * may be populated with additional properties.
	 */
}

// MicrosoftGraphPackage "@odata.type": "microsoft.graph.package"
type MicrosoftGraphPackage struct {
	Type string `json:"type"`
}

// MicrosoftGraphItemReference "@odata.type": "microsoft.graph.itemReference"
type MicrosoftGraphItemReference struct {
	DriveID       string                       `json:"driveId"`
	DriveType     string                       `json:"driveType"` // personal, business, documentLibrary
	ID            string                       `json:"id"`
	ListID        string                       `json:"listId"`
	Name          string                       `json:"name"`
	Path          string                       `json:"path"`
	ShareID       string                       `json:"shareId"`
	SharepointIDs *MicrosoftGraphSharepointIDs `json:"sharepointIds"`
	SiteID        string                       `json:"siteId"`
}

// MicrosoftGraphSharepointIDs "@odata.type": "microsoft.graph.sharepointIds"
type MicrosoftGraphSharepointIDs struct {
	ListID           string `json:"listId"`
	ListItemID       string `json:"listItemId"`
	ListItemUniqueID string `json:"listItemUniqueId"`
	SiteID           string `json:"siteId"`
	SiteURL          string `json:"siteUrl"`
	TenantID         string `json:"tenantId"`
	WebID            string `json:"webId"`
}

// MicrosoftGraphPhoto "@odata.type": "microsoft.graph.photo"
type MicrosoftGraphPhoto struct {
	CameraMake          string     `json:"cameraMake"`
	CameraModel         string     `json:"cameraModel"`
	ExposureDenominator float64    `json:"exposureDenominator"`
	ExposureNumerator   float64    `json:"exposureNumerator"`
	FNumber             float64    `json:"fNumber"`
	FocalLength         float64    `json:"focalLength"`
	ISO                 int32      `json:"iso"`
	TakenDateTime       *time.Time `json:"takenDateTime"`
}

// MicrosoftGraphPublicationFacet "@odata.type": "microsoft.graph.publicationFacet"
type MicrosoftGraphPublicationFacet struct {
	Level     string `json:"level"` // published, checkout
	VersionID string `json:"versionId"`
}

// MicrosoftGraphRemoteItem "@odata.type": "microsoft.graph.remoteItem"
type MicrosoftGraphRemoteItem struct {
	ID                   string                        `json:"id"` // identifier
	CreatedBy            *MicrosoftGraphIdentitySet    `json:"createdBy"`
	CreatedDateTime      *time.Time                    `json:"createdDateTime"`
	File                 *MicrosoftGraphFile           `json:"file"`
	FileSystemInfo       *MicrosoftGraphFileSystemInfo `json:"fileSystemInfo"`
	Folder               *MicrosoftGraphFolder         `json:"folder"`
	LastModifiedBy       *MicrosoftGraphIdentitySet    `json:"lastModifiedBy"`
	LastModifiedDateTime *time.Time                    `json:"lastModifiedDateTime"`
	Name                 string                        `json:"name"`
	Package              *MicrosoftGraphPackage        `json:"package"`
	ParentReference      *MicrosoftGraphItemReference  `json:"parentReference"`
	Shared               *MicrosoftGraphShared         `json:"shared"`
	SharepointIDs        *MicrosoftGraphSharepointIDs  `json:"sharepointIds"`
	SpecialFolder        *MicrosoftGraphSpecialFolder  `json:"specialFolder"`
	Size                 int64                         `json:"size"`
	WebDavURL            string                        `json:"webDavUrl"`
	WebURL               string                        `json:"webUrl"`
}

// MicrosoftGraphRoot "@odata.type": "microsoft.graph.root"
type MicrosoftGraphRoot struct {
	/* The Root facet indicates that an object is the top-most one in its hierarchy.
	 * The presence (non-null) of the facet value indicates that the object is the
	 * root. A null (or missing) value indicates the object is not the root.
	 *
	 * Note: While this facet is empty today, in future API revisions the facet may
	 * be populated with additional properties.
	 */
}

// MicrosoftGraphSearchResult "@odata.type": "microsoft.graph.searchResult"
type MicrosoftGraphSearchResult struct {
	OnClickTelemetryURL string `json:"onClickTelemetryUrl"`
}

// MicrosoftGraphShared "@odata.type": "microsoft.graph.shared"
type MicrosoftGraphShared struct {
	Owner          *MicrosoftGraphIdentitySet `json:"owner"`
	Scope          string                     `json:"scope"` // anonymous, organization, users
	SharedBy       *MicrosoftGraphIdentitySet `json:"sharedBy"`
	SharedDateTime *time.Time                 `json:"sharedDateTime"`
}

// MicrosoftGraphSpecialFolder "@odata.type": "microsoft.graph.specialFolder"
type MicrosoftGraphSpecialFolder struct {
	Name string `json:"name"`
}

// MicrosoftGraphVideo "@odata.type": "microsoft.graph.video"
type MicrosoftGraphVideo struct {
	AudioBitsPerSample    int32   `json:"audioBitsPerSample"`
	AudioChannels         int32   `json:"audioChannels"`
	AudioFormat           string  `json:"audioFormat"`
	AudioSamplesPerSecond int32   `json:"audioSamplesPerSecond"`
	Bitrate               int32   `json:"bitrate"`
	Duration              int32   `json:"duration"`
	FourCC                string  `json:"fourCC"`
	FrameRate             float64 `json:"frameRate"`
	Height                int32   `json:"height"`
	Width                 int32   `json:"width"`
}

// MicrosoftGraphItemActivity "@odata.type": "microsoft.graph.itemActivity"
type MicrosoftGraphItemActivity struct {
	ID     string                             `json:"id"` // identifier
	Action *MicrosoftGraphItemActionSet       `json:"action"`
	Actor  *MicrosoftGraphIdentitySet         `json:"actor"`
	Times  *MicrosoftGraphitemActivityTimeSet `json:"times"`

	/* relationships */
	DriveItem *MicrosoftGraphDriveItem `json:"driveItem"`
	ListItem  *MicrosoftGraphListItem  `json:"listItem"`
}

// MicrosoftGraphItemActionSet "@odata.type": "microsoft.graph.itemActionSet"
type MicrosoftGraphItemActionSet struct {
	Comment *MicrosoftGraphCommentAction `json:"comment"`
	Create  *MicrosoftGraphCreateAction  `json:"create"`
	Delete  *MicrosoftGraphDeleteAction  `json:"delete"`
	Edit    *MicrosoftGraphEditAction    `json:"edit"`
	Mention *MicrosoftGraphMentionAction `json:"mention"`
	Move    *MicrosoftGraphMoveAction    `json:"move"`
	Rename  *MicrosoftGraphRenameAction  `json:"rename"`
	Restore *MicrosoftGraphRestoreAction `json:"restore"`
	Share   *MicrosoftGraphShareAction   `json:"share"`
	Version *MicrosoftGraphVersionAction `json:"version"`
}

// MicrosoftGraphCommentAction "@odata.type": "microsoft.graph.commentAction"
type MicrosoftGraphCommentAction struct {
	IsReply      bool                        `json:"isReply"`
	ParentAuthor *MicrosoftGraphIdentitySet  `json:"parentAuthor"`
	Participants []MicrosoftGraphIdentitySet `json:"participants"`
}

// MicrosoftGraphCreateAction "@odata.type": "microsoft.graph.createAction"
type MicrosoftGraphCreateAction struct {
	/* The presence of the CreateAction resource on an itemActivity indicates that
	 * the activity created an item.
	 *
	 * Note: While this resource is empty today, in future API revisions it may be
	 * populated with additional properties.
	 */
}

// MicrosoftGraphDeleteAction "@odata.type": "microsoft.graph.deleteAction"
type MicrosoftGraphDeleteAction struct {
	Name       string `json:"name"`
	ObjectType string `json:"objectType"` // File, Folder
}

// MicrosoftGraphEditAction "@odata.type": "microsoft.graph.editAction"
type MicrosoftGraphEditAction struct {
	/* The presence of the EditAction resource on an itemActivity indicates that
	 * the activity edited an item.
	 *
	 * Note: While this resource is empty today, in future API revisions it may
	 * be populated with additional properties.
	 */
}

// MicrosoftGraphMentionAction "@odata.type": "microsoft.graph.mentionAction"
type MicrosoftGraphMentionAction struct {
	Mentionees []MicrosoftGraphIdentitySet `json:"mentionees"`
}

// MicrosoftGraphMoveAction "@odata.type": "microsoft.graph.moveAction"
type MicrosoftGraphMoveAction struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// MicrosoftGraphRenameAction "@odata.type": "microsoft.graph.renameAction"
type MicrosoftGraphRenameAction struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

// MicrosoftGraphRestoreAction "@odata.type": "microsoft.graph.restoreAction"
type MicrosoftGraphRestoreAction struct {
	/* The presence of the RestoreAction resource on an itemActivity indicates
	 * that the activity restored an item.
	 *
	 * Note: While this resource is empty today, in future API revisions it may
	 *  be populated with additional properties.
	 */
}

// MicrosoftGraphShareAction "@odata.type": "microsoft.graph.shareAction"
type MicrosoftGraphShareAction struct {
	Recipients []MicrosoftGraphIdentitySet `json:"recipients"`
}

// MicrosoftGraphVersionAction "@odata.type": "microsoft.graph.versionAction"
type MicrosoftGraphVersionAction struct {
	NewVersion string `json:"newVersion"`
}

// MicrosoftGraphListItem "@odata.type": "microsoft.graph.listItem"
type MicrosoftGraphListItem struct {
	ContentType   *MicrosoftGraphContentTypeInfo `json:"contentType"`
	Fields        *MicrosoftGraphFieldValueSet   `json:"fields"`
	SharepointIDs *MicrosoftGraphSharepointIDs   `json:"sharepointIds"`

	/* relationships */
	Activities []MicrosoftGraphItemActivity    `json:"activities"`
	DriveItem  *MicrosoftGraphDriveItem        `json:"driveItem"`
	Versions   []MicrosoftGraphListItemVersion `json:"versions"`

	/* inherited from baseItem */
	ID                   string                       `json:"id"`
	CreatedBy            *MicrosoftGraphIdentitySet   `json:"createdBy"`
	CreatedDateTime      *time.Time                   `json:"createdDateTime"`
	Description          string                       `json:"description"`
	ETag                 string                       `json:"eTag"`
	LastModifiedBy       *MicrosoftGraphIdentitySet   `json:"lastModifiedBy"`
	LastModifiedDateTime *time.Time                   `json:"lastModifiedDateTime"`
	Name                 string                       `json:"name"`
	ParentReference      *MicrosoftGraphItemReference `json:"parentReference"`
	WebURL               string                       `json:"webUrl"`
}

// MicrosoftGraphContentTypeInfo "@odata.type": "microsoft.graph.contentTypeInfo"
type MicrosoftGraphContentTypeInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MicrosoftGraphFieldValueSet "@odata.type": "microsoft.graph.fieldValueSet"
type MicrosoftGraphFieldValueSet struct {
	Author         string `json:"Author"`
	AuthorLookupID string `json:"AuthorLookupId"`
	Name           string `json:"Name"`
	Color          string `json:"Color"`
	Quantity       int32  `json:"Quantity"`
}

// MicrosoftGraphListItemVersion "@odata.type": "microsoft.graph.listItemVersion"
type MicrosoftGraphListItemVersion struct {
	Fields               *MicrosoftGraphFieldValueSet    `json:"fields"`
	ID                   string                          `json:"id"`
	LastModifiedBy       *MicrosoftGraphIdentitySet      `json:"lastModifiedBy"`
	LastModifiedDateTime *time.Time                      `json:"lastModifiedDateTime"`
	Published            *MicrosoftGraphPublicationFacet `json:"published"`
}

// MicrosoftGraphitemActivityTimeSet "@odata.type": "microsoft.graph.itemActivityTimeSet"
type MicrosoftGraphitemActivityTimeSet struct {
	ObservedDateTime *time.Time `json:"observedDateTime"`
	RecordedDateTime *time.Time `json:"recordedDateTime"`
}

// MicrosoftGraphPermission "@odata.type": "microsoft.graph.permission"
type MicrosoftGraphPermission struct {
	ID                  string                           `json:"id"`
	GrantedTo           *MicrosoftGraphIdentitySet       `json:"grantedTo"`
	GrantedToIdentities []MicrosoftGraphIdentitySet      `json:"grantedToIdentities"`
	InheritedFrom       *MicrosoftGraphItemReference     `json:"inheritedFrom"`
	Invitation          *MicrosoftGraphSharingInvitation `json:"invitation"`
	Link                *MicrosoftGraphSharingLink       `json:"link"`
	Roles               []string                         `json:"roles"` // read, write, sp.owner, sp.member
	ShareID             string                           `json:"shareId"`
}

// MicrosoftGraphSharingInvitation "@odata.type": "microsoft.graph.sharingInvitation"
type MicrosoftGraphSharingInvitation struct {
	Email          string                     `json:"email"`
	InvitedBy      *MicrosoftGraphIdentitySet `json:"invitedBy"`
	SignInRequired bool                       `json:"signInRequired"`
}

// MicrosoftGraphSharingLink "@odata.type": "microsoft.graph.sharingLink"
type MicrosoftGraphSharingLink struct {
	Application *MicrosoftGraphIdentity `json:"application"`
	Type        string                  `json:"type"`  // view, edit, embed
	Scope       string                  `json:"scope"` // anonymous, organization
	WebHTML     string                  `json:"webHtml"`
	WebURL      string                  `json:"webUrl"`
}

// MicrosoftGraphThumbnailSet "@odata.type": "microsoft.graph.thumbnailSet"
type MicrosoftGraphThumbnailSet struct {
	ID     string                   `json:"id"` // identifier
	Large  *MicrosoftGraphThumbnail `json:"large"`
	Medium *MicrosoftGraphThumbnail `json:"medium"`
	Small  *MicrosoftGraphThumbnail `json:"small"`
	Source *MicrosoftGraphThumbnail `json:"source"`
}

// MicrosoftGraphThumbnail "@odata.type": "microsoft.graph.thumbnail"
type MicrosoftGraphThumbnail struct {
	Content      *EdmDotStream `json:"content"`
	Height       int32         `json:"height"`
	SourceItemID string        `json:"sourceItemId"`
	URL          string        `json:"url"`
	Width        int32         `json:"width"`
}

// MicrosoftGraphDriveItemVersion "@odata.type": "microsoft.graph.driveItemVersion"
type MicrosoftGraphDriveItemVersion struct {
	Content              *EdmDotStream                   `json:"content"`
	ID                   string                          `json:"id"`
	LastModifiedBy       *MicrosoftGraphIdentitySet      `json:"lastModifiedBy"`
	LastModifiedDateTime *time.Time                      `json:"lastModifiedDateTime"`
	Publication          *MicrosoftGraphPublicationFacet `json:"publication"`
	Height               int32                           `json:"height"`
}

// MicrosoftGraphDrive  "@odata.type": "microsoft.graph.drive"
type MicrosoftGraphDrive struct {
	Activities           []MicrosoftGraphItemActivity `json:"activities"`
	ID                   string                       `json:"id"` // identifier
	CreatedBy            *MicrosoftGraphIdentitySet   `json:"createdBy"`
	CreatedDateTime      *time.Time                   `json:"createdDateTime"`
	Description          string                       `json:"description"`
	DriveType            string                       `json:"driveType"` // personal, business, documentLibrary
	Items                []MicrosoftGraphDriveItem    `json:"items"`
	LastModifiedBy       *MicrosoftGraphIdentitySet   `json:"lastModifiedBy"`
	LastModifiedDateTime *time.Time                   `json:"lastModifiedDateTime"`
	Name                 string                       `json:"name"`
	Owner                *MicrosoftGraphIdentitySet   `json:"owner"`
	Quota                *MicrosoftGraphQuota         `json:"quota"`
	Root                 *MicrosoftGraphDriveItem     `json:"root"`
	SharepointIDs        *MicrosoftGraphSharepointIDs `json:"sharepointIds"`
	Special              []MicrosoftGraphDriveItem    `json:"special"`
	System               *MicrosoftGraphSystemFacet   `json:"system"`
	WebURL               string                       `json:"webUrl"`
}

// MicrosoftGraphQuota "@odata.type": "microsoft.graph.quota"
type MicrosoftGraphQuota struct {
	Deleted   int64  `json:"deleted"`
	FileCount int64  `json:"fileCount"`
	Remaining int64  `json:"remaining"`
	State     string `json:"state"` // normal, nearing, critical, exceeded
	Total     int64  `json:"total"`
	Used      int64  `json:"used"`
}

// MicrosoftGraphSystemFacet "@odata.type": "microsoft.graph.systemFacet"
type MicrosoftGraphSystemFacet struct {
	/* The System facet indicates that the object is managed by the system
	 * for its own operation. Most apps should ignore items that have a
	 * System facet.
	 *
	 * Note: While this facet is empty today, in future API revisions the
	 * facet may be populated with additional properties.
	 */
}
