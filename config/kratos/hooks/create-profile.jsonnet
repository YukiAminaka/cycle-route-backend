function(ctx) {
  kratos_id: ctx.identity.id,
  email: if std.objectHas(ctx.identity.traits, 'email') then ctx.identity.traits.email else null,
  name: if std.objectHas(ctx.identity.traits, 'username') then ctx.identity.traits.username else ctx.identity.traits.email,
  first_name: if std.objectHas(ctx.identity.traits, 'name') && std.objectHas(ctx.identity.traits.name, 'first') then ctx.identity.traits.name.first else null,
  last_name: if std.objectHas(ctx.identity.traits, 'name') && std.objectHas(ctx.identity.traits.name, 'last') then ctx.identity.traits.name.last else null,
}
