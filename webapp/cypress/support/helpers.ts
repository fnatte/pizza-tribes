import { Education } from "../../src/generated/education";
import { Mouse } from "../../src/generated/gamestate";
import { generateId } from "../../src/utils";

export class MiceBuilder {
  private mice: Record<string, Partial<Mouse>> = {};


  add(education?: Education, count: number = 1): MiceBuilder {
    for (let n = 0; n < count; n++) {
      const id = generateId()
      this.mice[id] = {
        education,
        isEducated: education !== undefined,
        name:id,
      }
    }

    return this;
  }

  build(): Record<string, Partial<Mouse>> {
    return this.mice;
  }
}
