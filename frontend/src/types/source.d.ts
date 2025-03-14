type SourceModel = {
  id: string;
  name: string;
  type: string;
  settings: {[key: string]: any}
  auth_type: string;
  auth: {[key: string]: any}
}

type SourceModelChanges = {
  name: boolean;
  type: boolean;
  settings: boolean;
  auth: boolean;
}