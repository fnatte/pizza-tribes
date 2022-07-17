import type { PartialMessage, IMessageType } from "@protobuf-ts/runtime";
import { MESSAGE_TYPE, containsMessageType } from "@protobuf-ts/runtime";
import { get, set, setWith } from "lodash";
import { FieldMask } from "./generated/google/protobuf/field_mask";

export function extractMessage<T extends object>(msg: T, fieldMask: FieldMask) {
  const target: PartialMessage<T> = {};
  fieldMask.paths.forEach((x) =>
    extractField(target, msg, parseFieldMaskPath(x))
  );
  return target;
}

function parseFieldMaskPath(path: string): string[] {
  return path.split(".");
}

function extractField<T extends object>(
  target: PartialMessage<T>,
  msg: T,
  path: string[]
): PartialMessage<T> {
  setWith(target, path, get(msg, path), Object);

  /*
  let type: IMessageType<any>;
  if (containsMessageType(msg)) {
    type = msg[MESSAGE_TYPE];
  } else {
    throw new Error("Could not find message type on message.");
  }

  let curType: IMessageType<any> = type;
  let curMsg: any = msg;

  for (let i = 0; i < path.length; i++) {
    const field = curType.fields.find((x) => x.name === path[i]);
    if (field === undefined) {
      throw new Error("Could not find field: " + path);
    }

    if (field.kind === "message") {
      curMsg = curMsg[field.localName];
      curType = field.T();
    } else if (field.kind === "map") {
    } else if (field.kind === "scalar") {
    }

    if (i === path.length - 1) {
      if (field.kind === "scalar") {
      }
    }
  }
  */

  return target;
}
