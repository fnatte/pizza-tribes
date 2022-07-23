import type {
  PartialMessage,
  IMessageType,
  FieldInfo,
} from "@protobuf-ts/runtime";
import { MESSAGE_TYPE, containsMessageType } from "@protobuf-ts/runtime";
import { get, setWith } from "lodash";

export function extractMessage<T extends object>(msg: T, paths: string[]) {
  const target: PartialMessage<T> = {};
  paths.forEach((x) => extractField(target, msg, parseFieldMaskPath(x)));
  return target;
}

function parseFieldMaskPath(path: string): string[] {
  return path.split(".");
}

type FieldInfoMap = Extract<FieldInfo, { kind: 'map' }>;
type FieldInfoMapValue = FieldInfoMap['V'];

function findField(msg: object, path: string[]): FieldInfo | FieldInfoMapValue {
  let rootType: IMessageType<any>;
  if (containsMessageType(msg)) {
    rootType = msg[MESSAGE_TYPE];
  } else {
    throw new Error("Could not find message type on message.");
  }

  const field = rootType.fields.find((x) => x.localName == path[0]);
  if (!field) {
    throw new Error(`Failed to find field ${path[0]}`);
  }
  let curField: FieldInfo | FieldInfoMapValue = field;


  for (let i = 1; i < path.length; i++) {
    switch (curField.kind) {
      case "message":
        const field = curField.T().fields.find((x) => x.name == path[i]);
        if (!field) {
          throw new Error(`Failed to find field ${path[i]}`);
        }
        curField = field;
        break;
      case "map":
        curField = curField.V;
        break;
      case "enum":
      case "scalar":
        throw new Error(
          `${curField.kind} cannot have a nested field ${path[i]}`
        );
    }
  }

  return curField;
}

function extractField<T extends object>(
  target: PartialMessage<T>,
  msg: T,
  path: string[]
): PartialMessage<T> {
  let i = 0;
  const fields = path.map((_, i) => findField(msg, path.slice(0, i + 1)));
  const jsonPath = fields.map((f, i) => 'jsonName' in f ? f.jsonName : path[i]);

  for (let i = 0; i < path.length - 1; i++) {
    const field = fields[i];
    if ('repeat' in field && field.repeat) {
      throw new Error("A repeated field is only allowed at last path position");
    }
  }

  setWith(target, jsonPath, get(msg, jsonPath), (nsValue) => {
    if (nsValue) {
      return nsValue;
    }

    const field = fields[i];
    i++;
    return 'repeat' in field && field.repeat ? Array() : Object();
  });

  return target;
}
