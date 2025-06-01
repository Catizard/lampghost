<script lang="ts">
import { dto } from '@wailsjs/go/models';
import { NTag } from 'naive-ui';
import { defineComponent, h, VNode } from 'vue';

// NOTE: It's easier to write this component by using render function 
// instead of an explicit template block since the props passed to NaiveUI's
// NTag component are nullable
export default {
  props: {
    tableTags: Array<dto.DiffTableTagDto>,
  },
  setup(props) {
    return () => {
      const nodes: Array<VNode> = [];
      if (props.tableTags.length == 0) {
        return nodes;
      }
      props.tableTags.forEach(tag => {
        const props = {
          size: "small",
          style: {
            "margin-right": "5px"
          },
          color: {

          },
        };
        if (tag.TableTagColor.length > 0) {
          (props as any).color.color = tag.TableTagColor;
        }
        if (tag.TableTagTextColor.length > 0) {
          (props as any).color.textColor = tag.TableTagTextColor;
        }
        const node = h(
          NTag,
          props as any,
          { default: () => tag.TableSymbol + tag.TableLevel }
        )
        nodes.push(node);
      });
      return nodes;
    }
  }
}
</script>