export interface PublicRoom {
	id: string;
	options: Option[];
	question: string;
	ownedByMe: boolean;
}

export interface Room extends PublicRoom {
	ownerID: string;
}

export interface Option {
	id: string;
	ownerID: string;
	value: string;
}
