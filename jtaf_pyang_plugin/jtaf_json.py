# Copyright (c) 2021, Juniper Networks, Inc.
# All rights reserved.
"""jtaf output plugin

This plugin generates the JSON file used by jtaf It works on the entire
AST and ignores any pyang filters/paths passed on the command line.

"""

import optparse
import sys
import string
import pdb
import json
from collections import OrderedDict
import os
from pyang import plugin
from pyang import statements
from pyang import error
from pyang import types

sys.path.insert(0, os.path.dirname(__file__))


class FNode:
    """Represents a jtaf node in the configuration tree"""
    def __init__(self, name):
        self.name = name
        self.type = "container"
        self.children = []
        
    def __setitem__(self, a, v):
        """To allow setting of attributes directly
           eg. fn["my-attr"] = 5
        """
        self.__dict__[a] = v

    def __getitem__(self, a):
        return self.__dict__[a]
        
    def to_json_dict(self):
        """Orders the output of the json attributes for
           readability
        """
        keys = [
            "name",
            "type",
            "key",
            "leaf_type", 
            "other",
            "child",
        ]

        for k in self.__dict__:
            if not k in keys:
                keys.append(k)

        # make children last
        keys.append("children")
        
        d = OrderedDict()
        for k in keys:
            kk = k.replace("_", "-")
            if k in self.__dict__:
                # if children, output only if children is not empty
                if k != "children" or len(self.__dict__[k]) > 0:
                    d[kk] = self.__dict__[k]
                
        return d
        
class FNodeTree:
    def __init__(self):
        root = FNode("root")
        self.tree = {
            "root" : root,
            "identities" : []
        }
        self.stack= [self.tree["root"]]
        self.cur = root

    def to_json_dict(self):
        return self.tree
    
    def append_ident(self, ident):
        self.tree["identities"].append(ident)
    
    def push(self, name):
        child = FNode(name)
        if self.cur == None:
            self.tree["root"] = child
        else:
            self.cur.children.append(child)
            
        self.cur = child
        self.stack.append(child)
        
        return child

    def pop(self):
        if len(self.stack) > 0:
            cur =  self.stack[-1]
            self.stack = self.stack[:-1]
            if len(self.stack) > 0:
                self.cur = self.stack[-1]
            else:
                self.cur = None

    def top(self):
        if len(self.stack) > 0:
            return self.stack[-1]
        else:
            return None

    def set_attr(self, attr, val):
        n = self.top()
        if attr in yang_coerce:
            val = yang_coerce[attr](n, attr, val)


        n[attr] = val

    def add_prop(self, ch, key):
        for p in ch.search(key):
            if p != None:
                self.set_attr(p.keyword, p.arg)
 
#
# Prefix used in jtaf.yang module
#
fh_prefix = 'jtaf'

jt = FNodeTree()

yang_type = ["list", "container", "leaf", "leaf-list", "leafref", "choice", "case"]

# These are the yang properties that we output to json,
# all others are ignored

def to_bool(obj, key, val):
    if val == "true":
        return True
    else:
        return False

def to_int(obj, key, val):
    return int(val)

def to_array(obj, key, val):
    if hasattr(obj, key):
        return obj[key] + [val]
    else:
        return [val]


yang_coerce = {
    "mandatory" : to_bool,
    "min-elements" : to_int,
    "must" : to_array
    }

yang_props = ["key", "if-feature", "default", "must", "units", "mandatory",
              "min-elements", "ordered-by", "config", "path", "when"]

# These are explicitly skipped
yang_props_skip = ["type", "description"]


def pyang_plugin_init():
    plugin.register_plugin(FoghornPlugin())

class FoghornPlugin(plugin.PyangPlugin):
    def __init__(self):
        plugin.PyangPlugin.__init__(self, fh_prefix)
        
    def add_output_format(self, fmts):
        self.multiple_modules = True
        fmts[fh_prefix] = self

    def add_opts(self, optparser):
        optlist = [
            optparse.make_option("--jtaf-strip-apply",
                                 dest="jtaf_strip_apply",
                                 action="store_true",
                                 help="Don't output apply-* nodes"),
            optparse.make_option("--jtaf-no-restrictions",
                                 dest="jtaf_no_restrictions",
                                 action="store_true",
                                 help="Don't output restrictions"),
            ]
        g = optparser.add_option_group("jtaf output specific options")
        g.add_options(optlist)

    def setup_ctx(self, ctx):
        if ctx.opts.tree_help:
            jtaf_print_help()
            sys.exit(0)

    def setup_fmt(self, ctx):
        ctx.implicit_errors = False

    def emit(self, ctx, modules, fd):
        for epos, etag, eargs in ctx.errors:
            if error.is_error(error.err_level(etag)):
                raise error.EmitError("fhjson plugin needs a valid module (%s, %s, %s)" % (epos, etag, eargs))

        if ctx.opts.tree_path is not None:
            path = ctx.opts.tree_path.split('/')
            if path[0] == '':
                path = path[1:]
        else:
            path = None

        jtaf_emit_tree(ctx, fd, modules)


def jtaf_print_help():
    print("""
jtaf plugin
""")

def serialize(obj):
    if isinstance(obj, FNode):
        return obj.to_json_dict()
    if isinstance(obj, FNodeTree):
        return obj.to_json_dict()
    
    return obj.__dict__

#
# Generate the json from the statement-tree
#
def jtaf_emit_tree(ctx, fd, modules):
    for module in modules:
        jtaf_walk_identities(module)
        jtaf_walk_top_level(ctx, module)

    fd.write(json.dumps(jt, indent = 2, default=serialize))
    
def jtaf_get_keyvalue(stmt, kw):
    v = stmt.search((fh_prefix, kw))
    if len(v) == 1:
        return v[0].arg
    elif len(v) == 0:
        return None
    else:
        p = []
        for a in v:
            p.append(a)
        return a

def jtaf_walk_dts(node, ch):
    if ch.keyword[0] != fh_prefix:
        return
    node["dts_" + ch.keyword[1]] = ch.arg
    
def jtaf_walk_type(ctx, ch):
    t = ch.search_one("type")
    if t != None:
        leaf_type = t.arg
        if leaf_type == "identityref":
            b = t.search_one("base")
            if b != None:
                leaf_type = b.arg
        jt.set_attr("leaf-type", leaf_type)

        # Record built-in base type if this is typedef'd.
        base_type = get_base_type(t)
        if base_type.arg != leaf_type:
            jt.set_attr("base-type", base_type.arg)

        # Add any restrictions on the type, unless the user doesn't want us to
        # For enum, we always generate the restrictions, as jt depends on this.
        enum_type = (leaf_type == "enumeration") or (base_type.arg == "enumeration")
        if ctx.opts.jtaf_no_restrictions and not(enum_type):
            properties = {}
        else:
            properties = get_type_restrictions(t)
        for k, v in properties.items():
            jt.set_attr(k, v)

def jtaf_walk_identities(mod):
    if len(mod.i_identities.items()) == 0:
        return

    for identity in mod.i_identities.items():
        ch = identity[1]
        
        ident = {}
        ident["name"] = ch.arg
        for child in ch.substmts:
            if child.keyword in ['base', 'description']:
                ident[child.keyword] = child.arg
        jt.append_ident(ident)
    
def jtaf_walk_top_level(ctx, mod):
    if len(mod.i_children) == 0:
        return

    for ch in mod.i_children:
        if not ch.keyword in ["container", "list", "leaf", "leaf-list"]:
            continue

        jtaf_walk_child(ctx, ch)


def jtaf_walk_child(ctx, ch):
    if hasattr(ch, 'keyword') == False:
        return

    if hasattr(ch, 'arg') and ctx.opts.jtaf_strip_apply and ch.arg.startswith("apply-"):
        return
    
    if ch.keyword == 'uses':
        for child in ch.i_grouping.substmts:
            jtaf_walk_child(ctx, child)
        return

    if ch.keyword == 'typedef':
        jt.add_typedef(ch)
        return
    
    if type(ch.keyword) is tuple:
        jtaf_walk_dts(fh_cur_node, ch)
        return

    if ch.keyword in yang_type:
        jt.push(ch.arg)
        jt.set_attr("type", ch.keyword)
    else:
        jt.set_attr(ch.keyword, ch.arg)

    jtaf_walk_type(ctx, ch)
        
    for p in yang_props:
        jt.add_prop(ch, p)
        
    if hasattr(ch, 'i_children') and len(ch.i_children) > 0:
        for child in ch.i_children:
            if child.keyword in yang_props:
                continue

            if child.keyword in yang_props_skip:
                continue

            if not child.keyword in yang_type + ["uses"]:
                pass

            jtaf_walk_child(ctx, child)
    jt.pop()


def get_type_restrictions(type_stmt):
    """Return a dictionary of restrictions associated with the given type."""
    properties = {}
    add_restrictions(properties, type_stmt)
    return properties


def add_restrictions(parent, type_stmt):
    """
    Collect restrictions defined by a chain of type specs.

    Type statement defines the tail of the chain.
    Restrictions are added to the parent dictionary.
    """
    type_spec = type_stmt.i_type_spec

    while type_spec is not None:
        rest = fetch_restriction(type_spec)
        if rest is not None:
            attr_name = rest[0]
            attr_val = rest[1]
            if attr_name in parent:
                # Patterns from base types of a typedef'd need to be
                # included. Per RFC 7950 9.4.5:
                # "If a pattern restriction is applied to a type that is already
                #  pattern-restricted, values must match all patterns in the base
                #  type, in addition to the new patterns."
                # Lengths and ranges from base types must NOT be included, however.
                # Pyang will have rejected any length/range restrictions that are
                # less restrictive than those in a base type. It follows that only
                # the length/range restrictions in the derived type need to be
                # enforced.
                if attr_name == "patterns":
                    parent[attr_name] = parent[attr_name] + attr_val
            else:
                parent[attr_name] = attr_val
        type_spec = type_spec.base


def fetch_restriction(type_spec):
    """Return the restrictions from the specified type spec."""
    type_spec_name = type(type_spec).__name__
    if type_spec_name in type_specs:
        return type_specs[type_spec_name](type_spec)
    return None


def fetch_decimal64_restrictions(type_spec):
    """Return the restrictions from the Decimal64TypeSpec."""
    return "fraction-digits", type_spec.fraction_digits

def fetch_pattern_restrictions(type_spec):
    """Return the restrictions from the PatternTypeSpec."""
    patterns = []
    for pattern in type_spec.res:
        patterns.append(pattern.__str__())
    return "patterns", patterns


def fetch_length_restrictions(type_spec):
    """Return the restrictions from the LengthTypeSpec."""
    # Emit explicit min/max values.
    # The min/max keywords are replaced with the minimum/maximum values from the
    # type_spec (which will reflect the length restrictions defined by any
    # base types.
    # If no maximum length is specified, it defaults to the minimum value.
    lslist = []
    for lspec in type_spec.lengths:
        minlen = lspec[0] if lspec != "min" else type_spec.min
        if lspec[1] is not None:
            maxlen = lspec[1] if lspec[1] != "max" else type_spec.max
        else:
            maxlen = type_spec.min
        lslist.append({"min": minlen, "max": maxlen})
    return "lengths", lslist


def fetch_range_restrictions(type_spec):
    """Return the restrictions from the RangeTypeSpec."""
    ranges = []
    if type_spec.name == "decimal64":
        fetch_dec64_range_restrictions(ranges, type_spec)
    else:
        fetch_integer_range_restrictions(ranges, type_spec)

    return ("ranges", ranges) if len(ranges) > 0 else None


def fetch_integer_range_restrictions(ranges, type_spec):
    """Return the restrictions from the RangeTypeSpec for an integer."""
    # Emit explicit min/max values, but only if at least one of them is not
    # the default value for the type of integer.
    # The min/max keywords are replaced with the minimum/maximum values from the
    # type_spec (which will reflect any range restrictions defined by any base
    # types).
    # If no maximum length is specified, it defaults to the minimum value.
    default_spec = types.yang_type_specs[type_spec.name]
    for rdef in type_spec.ranges:
        if (rdef[0] != default_spec.min) or (rdef[1] != default_spec.max):
            rspec = {}
            rspec["min"] = rdef[0] if rdef[0] != "min" else type_spec.min
            if rdef[1] is not None:
                rspec["max"] = (
                    rdef[1] if rdef[1] != "max" else type_spec.max
                )
            else:
                rspec["max"] = rspec["min"]

            ranges.append(rspec)


def fetch_dec64_range_restrictions(ranges, type_spec):
    """Return the restrictions from the RangeTypeSpec for a decimal64."""
    # Emit explicit min/max values, allowing for values being defined either by a
    # Decimal64Value or by a scalar type.
    # The min/max keywords are replaced with the minimum/maximum values from the
    # type_spec (which will reflect any range restrictions defined by any
    # base types).
    # If no maximum length is specified, it defaults to the minimum value.
    for rdef in type_spec.ranges:
        rspec = {}
        if rdef[0] == "min":
            rspec["min"] = type_spec.min.s
        else:
            rspec["min"] = (
                rdef[0].value
                if isinstance(rdef[0], types.Decimal64Value)
                else rdef[0]
            )
        if rdef[1] is not None:
            if rdef[1] == "max":
                rspec["max"] = type_spec.max.s
            else:
                rspec["max"] = (
                    rdef[1].value
                    if isinstance(rdef[1], types.Decimal64Value)
                    else rdef[1]
                )
        else:
            rspec["max"] = rspec["min"]

        ranges.append(rspec)


def fetch_identityref_restrictions(type_spec):
    """Return the restrictions from the IdentityRefTypeSpec."""
    bases = []
    for base in type_spec.idbases:
        bases.append(base.arg)
    return "bases", bases


def fetch_path_restrictions(type_spec):
    """Return the restrictions from the PathTypeSpec."""
    return "path", type_spec.path_.arg


def fetch_inst_identifier_restrictions(type_spec):
    """Return the restrictions from the InstanceIdentifierTypeSpec."""
    return "require-instance", type_spec.require_instance


def fetch_leafref_restrictions(type_spec):
    """Return the restrictions from the LeafrefTypeSpec."""
    return "require-instance", type_spec.require_instance


def fetch_enum_restrictions(type_spec):
    """Return the restrictions from the EnumTypeSpec."""
    enums = []
    for enum in type_spec.enums:
        enums.append({"id": enum[0], "value": enum[1]})
    return "enums", enums


def fetch_bit_restrictions(type_spec):
    """Return the restrictions from the BitTypeSpec."""
    return "bits", type_spec.bits


def fetch_union_restrictions(type_spec):
    """Return the restrictions from the UnionTypeSpec."""
    utypes = []
    for type_stmt in type_spec.types:
        base_type = get_base_type(type_stmt)
        utype = {"type": base_type.arg}
        properties = get_type_restrictions(type_stmt)
        utype.update(properties)
        utypes.append(utype)

    return "types", utypes

# Associate the type of pyang type_spec with the method used to extract
# restrictions from it.
type_specs = {
        "PatternTypeSpec": fetch_pattern_restrictions,
        "EnumTypeSpec": fetch_enum_restrictions,
        "IdentityrefTypeSpec": fetch_identityref_restrictions,
        "Decimal64TypeSpec": fetch_decimal64_restrictions,
        "PathTypeSpec": fetch_path_restrictions,
        "RangeTypeSpec": fetch_range_restrictions,
        "LengthTypeSpec": fetch_length_restrictions,
        "BitTypeSpec": fetch_bit_restrictions,
        "LeafrefTypeSpec": fetch_leafref_restrictions,
        "InstanceIdentifierTypeSpec": fetch_inst_identifier_restrictions,
        "UnionTypeSpec": fetch_union_restrictions,
    }

def get_base_type(stmt):
    """Return the built-in type from which the stmt is derived."""
    try:
        typedef = stmt.i_typedef

    except AttributeError:
        return stmt
    else:
        if typedef is not None:
            return get_base_type(typedef.search_one("type"))

    return stmt
