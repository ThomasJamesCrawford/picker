export interface PublicRoom {
	id: string;
	options: Option[];
	question: string;
}

export interface Room extends PublicRoom {
	ownerID: string;
}

export interface Option {
	id: string;
	ownerID: string;
	value: string;
}
